package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"
)

const MGMT_MSG_MARKER_NATIVE = 0x23232301

const (
	MSG_CODE_SESSION_REQ  = 10
	MSG_CODE_SESSION_REPLY = 11
	MSG_CODE_NOTIFY_SELECT = 9
	MSG_CODE_NOTIFY       = 4
)

const MSG_FORMAT_JSON = 2

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

		if err := readSessionReply(conn); err != nil {
			fmt.Printf("Failed to read session reply: %v\n", err)
			conn.Close()
			continue
		}
		fmt.Println("=== Session established ===")

		// Subscribe to notifications
		if err := sendNotifySelect(conn); err != nil {
			fmt.Printf("Failed to send notify select: %v\n", err)
			conn.Close()
			continue
		}
		fmt.Println("=== Subscribed to notifications ===")

		// Read messages loop
		buf := make([]byte, 65536)
		for {
			if _, err := io.ReadFull(conn, buf[:8]); err != nil {
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

			if bodyLen >= 24 {
				code := binary.LittleEndian.Uint16(buf[8:10])
				referID := binary.LittleEndian.Uint64(buf[14:22])
				reqID := binary.LittleEndian.Uint64(buf[22:30])

				if code == MSG_CODE_NOTIFY {
					fmt.Println("\n=== NOTIFICATION RECEIVED ===")
					fmt.Printf("Refer ID: %d, Req ID: %d\n", referID, reqID)
					if bodyLen > 32 {
						fmt.Printf("Data: %s\n", string(buf[32:bodyLen]))
					}
				} else {
					fmt.Printf("Message: code=%d refer=%d req=%d\n", code, referID, reqID)
				}
			}
		}

		conn.Close()
		time.Sleep(1 * time.Second)
	}
}

func sendSessionRequest(conn net.Conn) error {
	header := make([]byte, 24)
	binary.LittleEndian.PutUint16(header[0:2], MSG_CODE_SESSION_REQ)
	binary.LittleEndian.PutUint64(header[8:16], 0)
	binary.LittleEndian.PutUint64(header[16:24], 1)

	mdata := make([]byte, 24+1+7+11)
	copy(mdata[0:24], header)
	mdata[24] = MSG_FORMAT_JSON
	mdata[32] = 0
	copy(mdata[33:], []byte("go-watcher"))

	msgLen := uint32(8 + len(mdata))
	msg := make([]byte, msgLen)
	binary.LittleEndian.PutUint32(msg[0:4], MGMT_MSG_MARKER_NATIVE)
	binary.LittleEndian.PutUint32(msg[4:8], msgLen)
	copy(msg[8:], mdata)

	_, err := conn.Write(msg)
	return err
}

func readSessionReply(conn net.Conn) error {
	buf := make([]byte, 65536)
	if _, err := io.ReadFull(conn, buf[:8]); err != nil {
		return err
	}
	marker := binary.LittleEndian.Uint32(buf[:4])
	if marker != MGMT_MSG_MARKER_NATIVE {
		return fmt.Errorf("bad marker")
	}
	length := binary.LittleEndian.Uint32(buf[4:8])
	_, err := io.ReadFull(conn, buf[8:8+int(length)-8])
	return err
}

func sendNotifySelect(conn net.Conn) error {
	header := make([]byte, 24)
	binary.LittleEndian.PutUint16(header[0:2], MSG_CODE_NOTIFY_SELECT)
	binary.LittleEndian.PutUint64(header[8:16], 0)
	binary.LittleEndian.PutUint64(header[16:24], 1)

	body := make([]byte, 24+1+7+2)
	copy(body[0:24], header)
	body[24] = 1
	body[32] = 0
	copy(body[33:], []byte("/"))

	msgLen := 8 + len(body)
	msg := make([]byte, msgLen)
	binary.LittleEndian.PutUint32(msg[0:4], MGMT_MSG_MARKER_NATIVE)
	binary.LittleEndian.PutUint32(msg[4:8], msgLen)
	copy(msg[8:], body)

	_, err := conn.Write(msg)
	return err
}
