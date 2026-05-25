package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sort"
	"strings"
	"time"
)

const MGMT_MSG_MARKER_NATIVE = 0x23232301

const (
	MSG_CODE_ERROR         = 0
	MSG_CODE_TREE_DATA     = 2
	MSG_CODE_GET_DATA      = 3
	MSG_CODE_SESSION_REQ   = 10
	MSG_CODE_SESSION_REPLY = 11
	MSG_CODE_NOTIFY        = 4
)

const MSG_FORMAT_JSON = 2
const routeMapXPath = "/frr-route-map:lib"
const backendClientsXPath = "/frr-backend:clients/client"

const (
	getDataFlagState        = 0x01
	getDataFlagConfig       = 0x02
	getDataDefaultsExplicit = 0
	datastoreRunning        = 1
	datastoreOperational    = 3
)

const (
	notifyOpNotification = 0
	notifyOpReplace      = 1
	notifyOpDelete       = 2
	notifyOpPatch        = 3
	notifyOpSync         = 4
)

func main() {
	fmt.Println("=== FRR Notification Watcher ===")
	fmt.Println("Listening on /run/frr/mgmtd_fe.sock...")

	for {
		conn, err := net.Dial("unix", "/run/frr/mgmtd_fe.sock")
		if err != nil {
			fmt.Printf("Connection failed: %v, retrying in 1s...\n", err)
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Println("=== Connected to MGMTD ===")

		// Send session request and read reply
		if err := sendSessionRequest(conn); err != nil {
			fmt.Printf("Failed to send session request: %v\n", err)
			conn.Close()
			continue
		}
		fmt.Println("=== Sent session request ===")

		sessID, err := readSessionReply(conn)
		if err != nil {
			fmt.Printf("Failed to read session reply: %v\n", err)
			conn.Close()
			continue
		}
		fmt.Printf("=== Session established: %d ===\n", sessID)

		fmt.Println("=== Listening for all YANG notifications (no notify selectors) ===")
		fmt.Printf("=== Polling route-map config via MGMTD GET_DATA: %s ===\n", routeMapXPath)
		fmt.Printf("=== Polling backend state via MGMTD GET_DATA: %s ===\n", backendClientsXPath)
		fmt.Println("=== Waiting for notifications ===")

		var lastRouteMapKey string
		var lastBackendClientsKey string
		pendingRequests := map[uint64]string{}
		nextReqID := uint64(2)
		lastPoll := time.Time{}
		buf := make([]byte, 65536)
		for {
			if time.Since(lastPoll) >= 250*time.Millisecond {
				lastPoll = time.Now()
				if !hasPendingRequest(pendingRequests, routeMapXPath) {
					reqID := nextReqID
					nextReqID++
					if err := sendGetData(conn, sessID, reqID, routeMapXPath, datastoreRunning, getDataFlagConfig); err != nil {
						fmt.Printf("Route-map poll error: %v\n", err)
						break
					}
					pendingRequests[reqID] = routeMapXPath
				}
				if !hasPendingRequest(pendingRequests, backendClientsXPath) {
					reqID := nextReqID
					nextReqID++
					if err := sendGetData(conn, sessID, reqID, backendClientsXPath, datastoreOperational, getDataFlagState); err != nil {
						fmt.Printf("Backend poll error: %v\n", err)
						break
					}
					pendingRequests[reqID] = backendClientsXPath
				}
			}

			_ = conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
			if _, err := io.ReadFull(conn, buf[:8]); err != nil {
				if isTimeout(err) {
					continue
				}
				fmt.Printf("Read error: %v\n", err)
				break
			}

			marker := binary.LittleEndian.Uint32(buf[:4])
			if marker != MGMT_MSG_MARKER_NATIVE {
				fmt.Printf("Bad marker: 0x%X\n", marker)
				break
			}

			length := binary.LittleEndian.Uint32(buf[4:8])
			bodyLen := int(length) - 8
			if bodyLen > len(buf)-8 {
				fmt.Printf("Message too large: %d bytes\n", bodyLen)
				break
			}

			if _, err := io.ReadFull(conn, buf[8:8+bodyLen]); err != nil {
				fmt.Printf("Body read error: %v\n", err)
				break
			}

			if bodyLen >= 32 {
				code := binary.LittleEndian.Uint16(buf[8:10])
				vsplit := binary.LittleEndian.Uint32(buf[12:16])
				msgSessID := binary.LittleEndian.Uint64(buf[16:24])
				reqID := binary.LittleEndian.Uint64(buf[24:32])

				if code == MSG_CODE_NOTIFY {
					resultType := buf[32]
					op := buf[33]
					payload := buf[40 : 8+bodyLen]
					xpath, data := splitNotificationPayload(payload, int(vsplit))

					fmt.Println("\n=== NOTIFICATION RECEIVED ===")
					fmt.Printf("Session ID: %d, Req ID: %d, Result Type: %d, Op: %d\n", msgSessID, reqID, resultType, op)
					printNotification(op, xpath, data)
				} else if code == MSG_CODE_TREE_DATA {
					xpath := pendingRequests[reqID]
					delete(pendingRequests, reqID)
					data := ""
					if bodyLen > 32 {
						data = cstr(buf[40 : 8+bodyLen])
					}
					switch xpath {
					case routeMapXPath:
						currentKey := canonicalJSON(data)
						if currentKey != lastRouteMapKey {
							if lastRouteMapKey != "" || hasRouteMapData(data) {
								fmt.Println("\n=== ROUTE-MAP CONFIG CHANGE ===")
								if hasRouteMapData(data) {
									fmt.Printf("#OP=REPLACE: %s\n", routeMapXPath)
									fmt.Println(data)
								} else {
									fmt.Printf("#OP=DELETE: %s\n", routeMapXPath)
								}
							}
							lastRouteMapKey = currentKey
						}
					case backendClientsXPath:
						currentKey := canonicalBackendClientsJSON(data)
						if currentKey != lastBackendClientsKey {
							if lastBackendClientsKey != "" {
								fmt.Println("\n=== BACKEND STATE CHANGE ===")
								fmt.Printf("#OP=PATCH: %s\n", backendClientsXPath)
								fmt.Println(data)
							}
							lastBackendClientsKey = currentKey
						}
					default:
						fmt.Printf("TREE_DATA reply: session=%d req=%d\n", msgSessID, reqID)
						if data != "" {
							fmt.Println(data)
						}
					}
				} else if code == MSG_CODE_ERROR {
					errCode := int16(binary.LittleEndian.Uint16(buf[32:34]))
					errMsg := cstr(buf[40 : 8+bodyLen])
					delete(pendingRequests, reqID)
					fmt.Printf("MGMTD error: session=%d req=%d error=%d message=%s\n", msgSessID, reqID, errCode, errMsg)
				} else {
					fmt.Printf("Message: code=%d session=%d req=%d\n", code, msgSessID, reqID)
				}
			}
		}

		conn.Close()
		time.Sleep(1 * time.Second)
	}
}

func sendSessionRequest(conn net.Conn) error {
	clientName := []byte("go-watcher\x00")

	header := make([]byte, 24)
	binary.LittleEndian.PutUint16(header[0:2], MSG_CODE_SESSION_REQ)
	binary.LittleEndian.PutUint64(header[8:16], 0)
	binary.LittleEndian.PutUint64(header[16:24], 1)

	mdata := make([]byte, 24+8+len(clientName))
	copy(mdata[0:24], header)
	mdata[24] = MSG_FORMAT_JSON
	copy(mdata[32:], clientName)

	msgLen := uint32(8 + len(mdata))
	msg := make([]byte, msgLen)
	binary.LittleEndian.PutUint32(msg[0:4], MGMT_MSG_MARKER_NATIVE)
	binary.LittleEndian.PutUint32(msg[4:8], msgLen)
	copy(msg[8:], mdata)

	_, err := conn.Write(msg)
	return err
}

func readSessionReply(conn net.Conn) (uint64, error) {
	buf := make([]byte, 65536)
	if _, err := io.ReadFull(conn, buf[:8]); err != nil {
		return 0, err
	}
	marker := binary.LittleEndian.Uint32(buf[:4])
	if marker != MGMT_MSG_MARKER_NATIVE {
		return 0, fmt.Errorf("bad marker")
	}
	length := binary.LittleEndian.Uint32(buf[4:8])
	bodyLen := int(length) - 8
	if bodyLen < 32 {
		return 0, fmt.Errorf("short session reply: %d bytes", bodyLen)
	}
	if _, err := io.ReadFull(conn, buf[8:8+bodyLen]); err != nil {
		return 0, err
	}

	code := binary.LittleEndian.Uint16(buf[8:10])
	if code == MSG_CODE_ERROR {
		return 0, fmt.Errorf("mgmtd error: %s", cstr(buf[40:8+bodyLen]))
	}
	if code != MSG_CODE_SESSION_REPLY {
		return 0, fmt.Errorf("expected session reply, got code=%d", code)
	}
	if buf[32] == 0 {
		return 0, fmt.Errorf("mgmtd did not create session")
	}

	return binary.LittleEndian.Uint64(buf[16:24]), nil
}

func sendGetData(conn net.Conn, sessID uint64, reqID uint64, xpath string, datastore byte, flags byte) error {
	query := []byte(xpath + "\x00")

	header := make([]byte, 24)
	binary.LittleEndian.PutUint16(header[0:2], MSG_CODE_GET_DATA)
	binary.LittleEndian.PutUint64(header[8:16], sessID)
	binary.LittleEndian.PutUint64(header[16:24], reqID)

	body := make([]byte, 24+8+len(query))
	copy(body[0:24], header)
	body[24] = MSG_FORMAT_JSON
	body[25] = flags
	body[26] = getDataDefaultsExplicit
	body[27] = datastore
	copy(body[32:], query)

	msgLen := uint32(8 + len(body))
	msg := make([]byte, msgLen)
	binary.LittleEndian.PutUint32(msg[0:4], MGMT_MSG_MARKER_NATIVE)
	binary.LittleEndian.PutUint32(msg[4:8], msgLen)
	copy(msg[8:], body)

	_, err := conn.Write(msg)
	return err
}

func splitNotificationPayload(payload []byte, vsplit int) (string, string) {
	if len(payload) == 0 {
		return "", ""
	}

	xpathEnd := len(payload)
	dataStart := len(payload)
	if vsplit > 0 && vsplit <= len(payload) {
		xpathEnd = vsplit - 1
		dataStart = vsplit
	}

	xpath := cstr(payload[:xpathEnd])
	data := ""
	if dataStart < len(payload) {
		data = cstr(payload[dataStart:])
	}

	return xpath, data
}

func printNotification(op byte, xpath string, data string) {
	switch op {
	case notifyOpNotification:
		if data != "" {
			fmt.Println(data)
		}
	case notifyOpReplace:
		fmt.Printf("#OP=REPLACE: %s\n", xpath)
		fmt.Println(data)
	case notifyOpDelete:
		fmt.Printf("#OP=DELETE: %s\n", xpath)
	case notifyOpPatch:
		fmt.Printf("#OP=PATCH: %s\n", xpath)
		fmt.Println(data)
	case notifyOpSync:
		fmt.Printf("#OP=SYNC: %s\n", xpath)
		fmt.Println(data)
	default:
		fmt.Printf("#OP=UNKNOWN(%d): %s\n", op, xpath)
		if data != "" {
			fmt.Println(data)
		}
	}
}

func cstr(data []byte) string {
	return strings.TrimRight(string(data), "\x00")
}

func hasRouteMapData(data string) bool {
	return strings.Contains(data, `"frr-route-map:lib"`)
}

func isTimeout(err error) bool {
	netErr, ok := err.(net.Error)
	return ok && netErr.Timeout()
}

func hasPendingRequest(pending map[uint64]string, xpath string) bool {
	for _, pendingXPath := range pending {
		if pendingXPath == xpath {
			return true
		}
	}
	return false
}

func canonicalJSON(data string) string {
	var value any
	if err := json.Unmarshal([]byte(data), &value); err != nil {
		return data
	}
	encoded, err := json.Marshal(value)
	if err != nil {
		return data
	}
	return string(encoded)
}

func canonicalBackendClientsJSON(data string) string {
	var payload struct {
		Clients struct {
			Client []map[string]any `json:"client"`
		} `json:"frr-backend:clients"`
	}
	if err := json.Unmarshal([]byte(data), &payload); err != nil {
		return data
	}
	sort.Slice(payload.Clients.Client, func(i, j int) bool {
		left, _ := payload.Clients.Client[i]["name"].(string)
		right, _ := payload.Clients.Client[j]["name"].(string)
		return left < right
	})
	encoded, err := json.Marshal(payload)
	if err != nil {
		return data
	}
	return string(encoded)
}
