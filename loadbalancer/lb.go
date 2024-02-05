package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var backends = []string{"backend1:8080", "backend2:8080", "backend3:8080"}
var currentServer int
var mutex sync.Mutex
var (
	healthCheckURL    string
	healthCheckPeriod int
)
var (
	activeBackends []string
)

func checkBackendHealth(backend string) bool {
	resp, err := http.Get("http://" + backend + healthCheckURL)
	if err != nil || resp.StatusCode != 200 {
		return false
	}
	return true
}

func updateActiveBackends() {
	for {
		mutex.Lock()
		activeBackends = activeBackends[:0]
		for _, backend := range backends {
			if checkBackendHealth(backend) {
				activeBackends = append(activeBackends, backend)
			}
		}
		mutex.Unlock()
		time.Sleep(time.Duration(healthCheckPeriod) * time.Second)
	}
}

func getNextBackend() string {
	mutex.Lock()
	defer mutex.Unlock()
	if len(activeBackends) == 0 {
		return ""
	}
	backend := activeBackends[currentServer%len(activeBackends)]
	currentServer++
	return backend
}

func main() {
	// health check
	flag.StringVar(&healthCheckURL, "health-check-url", "/health", "URL for health checking")
	flag.IntVar(&healthCheckPeriod, "health-check-period", 10, "Health check period in seconds")
	flag.Parse()
	go updateActiveBackends()

	// starts to listen
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
