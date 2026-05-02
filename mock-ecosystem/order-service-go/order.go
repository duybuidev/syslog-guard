package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
)

func main() {
	log.SetOutput(os.Stdout) // Ghi log ra Stdout để Docker thu thập

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Order Service is healthy"))
	})

	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		orderID := rand.Intn(10000)
		log.Printf("[INFO] Created order #%d successfully\n", orderID)
		w.Write([]byte(fmt.Sprintf("Order %d created", orderID)))
	})

	// API này dùng để "bắn sập" container, test SysWatch
	http.HandleFunc("/crash", func(w http.ResponseWriter, r *http.Request) {
		log.Println("[FATAL] Out of memory simulation. Crashing...")
		os.Exit(1)
	})

	port := "8081"
	log.Printf("Order Service starting on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %s", err)
	}
}
