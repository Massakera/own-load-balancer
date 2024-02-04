package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
)

func main() {
	listenAddr := ":80"
	backendAddr := "http://backend:8080"

	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", listenAddr, err)
	}
	defer listener.Close()
	log.Printf("Listening on %s\n", listenAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()

			request, err := http.ReadRequest(bufio.NewReader(c))
			if err != nil {
				log.Printf("Failed to read request: %v", err)
				return
			}

			log.Printf("Received request from %s\n%s %s %s\nHost: %s\nUser-Agent: %s\nAccept: %s\n",
				c.RemoteAddr(), request.Method, request.URL, request.Proto, request.Host, request.UserAgent(), request.Header.Get("Accept"))

			backendReq, err := http.NewRequest(request.Method, backendAddr+request.URL.String(), nil)
			if err != nil {
				log.Printf("Failed to create request for backend: %v", err)
				return
			}

			backendReq.Header = request.Header

			if request.Body != nil {
				backendReq.Body = io.NopCloser(request.Body)
			}

			client := &http.Client{}
			resp, err := client.Do(backendReq)
			if err != nil {
				log.Printf("Failed to send request to backend: %v", err)
				return
			}
			defer resp.Body.Close()

			err = resp.Write(c)
			if err != nil {
				log.Printf("Failed to write response to client: %v", err)
				return
			}
		}(conn)
	}
}
