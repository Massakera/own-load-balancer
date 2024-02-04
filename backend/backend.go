package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		log.Printf("Received request from %s\n%s %s %s\nHost: %s\nUser-Agent: %s\nAccept: %s\n",
			r.RemoteAddr, r.Method, r.URL.Path, r.Proto, r.Host, r.UserAgent(), r.Header.Get("Accept"))

		w.WriteHeader(http.StatusOK) // This sets the HTTP status to 200 OK
		fmt.Fprintf(w, "Hello From Backend Server")
	})

	log.Println("Listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to listen and serve: %v", err)
	}
}
