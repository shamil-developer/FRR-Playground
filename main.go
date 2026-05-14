package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	pb "frr-playground/frrgrpc"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	frrSocketPath = "/var/run/frr/bgpd.vty"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", healthHandler)

	mux.HandleFunc(
		"/frr/vtysh/running-config",
		vtyshRunningConfigHandler,
	)

	mux.HandleFunc(
		"/frr/socket/running-config",
		socketRunningConfigHandler,
	)

	mux.HandleFunc(
		"/frr/grpc/check",
		grpcCheckHandler,
	)

	mux.HandleFunc(
		"/frr/grpc/create-candidate",
		grpcCreateCandidateHandler,
	)

	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}

	go func() {
		log.Println("server started on :8081")

		if err := server.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped")
}

func healthHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusOK)

	_, _ = w.Write([]byte("ok"))
}

func vtyshRunningConfigHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, cancel := context.WithTimeout(
		r.Context(),
		5*time.Second,
	)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"vtysh",
		"-c",
		"show running-config",
	)

	out, err := cmd.CombinedOutput()

	if err != nil {
		writeResponse(w, nil, err)
		return
	}

	w.Header().Set(
		"Content-Type",
		"text/plain",
	)

	w.WriteHeader(http.StatusOK)

	_, _ = w.Write(out)
}

func socketRunningConfigHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, cancel := context.WithTimeout(
		r.Context(),
		5*time.Second,
	)
	defer cancel()

	out, err := readViaSocket(
		ctx,
		"show bgp summary json",
	)

	if err != nil {
		writeResponse(w, nil, err)
		return
	}

	writeResponse(w, out, nil)
}

func grpcCheckHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, cancel := context.WithTimeout(
		r.Context(),
		10*time.Second,
	)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		"localhost:50051",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithBlock(),
	)

	if err != nil {
		writeResponse(w, nil, err)
		return
	}

	defer conn.Close()

	client := pb.NewNorthboundClient(conn)

	req := &pb.GetRequest{
		Type:     pb.GetRequest_ALL,
		Encoding: pb.Encoding_JSON,
		Path: []string{
			"/frr-interface:lib",
		},
	}

	stream, err := client.Get(ctx, req)
	if err != nil {
		writeResponse(w, nil, err)
		return
	}

	var responses []any

	for {
		resp, err := stream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			writeResponse(w, nil, err)
			return
		}

		var parsed any

		err = json.Unmarshal(
			[]byte(resp.Data.Data),
			&parsed,
		)

		if err != nil {
			responses = append(
				responses,
				resp.Data.Data,
			)

			continue
		}

		responses = append(
			responses,
			parsed,
		)
	}

	writeResponse(w, responses, nil)
}

func readViaSocket(
	ctx context.Context,
	command string,
) (any, error) {

	conn, err := net.DialTimeout(
		"unix",
		frrSocketPath,
		3*time.Second,
	)

	if err != nil {
		return nil, err
	}

	defer conn.Close()

	err = conn.SetDeadline(
		time.Now().Add(3 * time.Second),
	)

	if err != nil {
		return nil, err
	}

	/*
		FRR protocol:
		command must end with NULL byte
	*/

	payload := append(
		[]byte(command),
		0x00,
	)

	_, err = conn.Write(payload)

	if err != nil {
		return nil, err
	}

	var result bytes.Buffer

	buf := make([]byte, 4096)

	for {

		n, err := conn.Read(buf)

		if n > 0 {

			result.Write(buf[:n])

			raw := result.Bytes()

			/*
				FRR response:

				payload\0\0\0status

				ищем marker
			*/

			if bytes.LastIndex(
				raw,
				[]byte{0x00, 0x00, 0x00},
			) != -1 {
				break
			}
		}

		if err != nil {

			/*
				FRR often keeps socket open.
				timeout here is NORMAL.
			*/

			if netErr, ok := err.(net.Error); ok &&
				netErr.Timeout() {
				break
			}

			if errors.Is(err, io.EOF) {
				break
			}

			return nil, err
		}
	}

	raw := result.Bytes()

	marker := bytes.LastIndex(
		raw,
		[]byte{0x00, 0x00, 0x00},
	)

	if marker == -1 {
		return nil, errors.New(
			"FRR marker not found",
		)
	}

	data := bytes.TrimRight(
		raw[:marker],
		"\x00",
	)

	if len(data) == 0 {
		return nil, errors.New(
			"empty FRR payload",
		)
	}

	var parsed any

	err = json.Unmarshal(
		data,
		&parsed,
	)

	if err != nil {
		return string(data), nil
	}

	return parsed, nil
}

func writeResponse(
	w http.ResponseWriter,
	result any,
	err error,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	if err != nil {
		w.WriteHeader(
			http.StatusInternalServerError,
		)

		_ = json.NewEncoder(w).Encode(
			map[string]string{
				"error": err.Error(),
			},
		)

		return
	}

	encoder := json.NewEncoder(w)

	encoder.SetIndent("", "  ")

	_ = encoder.Encode(result)
}

func grpcCreateCandidateHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, cancel := context.WithTimeout(
		r.Context(),
		10*time.Second,
	)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		"localhost:50051",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
		grpc.WithBlock(),
	)

	if err != nil {
		writeResponse(w, nil, err)
		return
	}

	defer conn.Close()

	client := pb.NewNorthboundClient(conn)

	/*
		Создаем candidate config
	*/

	resp, err := client.CreateCandidate(
		ctx,
		&pb.CreateCandidateRequest{},
	)

	if err != nil {
		writeResponse(w, nil, err)
		return
	}

	writeResponse(
		w,
		map[string]any{
			"candidateId": resp.CandidateId,
		},
		nil,
	)
}
