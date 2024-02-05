package main

import (
	"bufio"
	"log"
	"net"
	"net/http"
	"sync"
)

var backends = []string{"backend1:8080", "backend2:8081", "backend3:8082"}
var currentServer int
var mutex sync.Mutex

func getNextBackend() string {
	mutex.Lock()
	defer mutex.Unlock()
	backend := backends[currentServer]
	currentServer = (currentServer + 1) % len(backends)
	return backend
}

func main() {
	listenAddr := ":80"

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

			backendServer := getNextBackend()

			backendReq, err := http.NewRequest(request.Method, "http://"+backendServer+request.URL.String(), request.Body)
			if err != nil {
				log.Printf("Failed to create request for backend: %v", err)
				return
			}

			backendReq.Header = request.Header

			resp, err := http.DefaultClient.Do(backendReq)
			if err != nil {
				log.Printf("Failed to send request to backend: %v", err)
				return
			}
			defer resp.Body.Close()

			resp.Write(c)
		}(conn)
	}
}
