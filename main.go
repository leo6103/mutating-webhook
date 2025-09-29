package main

import (
	"fmt"
	"log"
	"net/http"
)

func mutateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "helloworld")
}

func main() {
	http.HandleFunc("/mutate", mutateHandler)

	addr := ":8080"
	log.Printf("🚀 Webhook server listening on %s\n", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
