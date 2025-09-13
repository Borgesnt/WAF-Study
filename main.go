package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

// LogEntry representa um log de requisição bloqueada
type LogEntry struct {
	Timestamp  string `json:"timestamp"`
	RemoteAddr string `json:"remote_addr"`
	Path       string `json:"path"`
	Reason     string `json:"reason"`
}

// Global state to store blocked logs
var (
	blockedLogs = make([]LogEntry, 0)
	mu          sync.Mutex
)

// WAFHandler é o nosso handler principal que atua como proxy e API
func WAFHandler(w http.ResponseWriter, r *http.Request) {
	// API Endpoints para o frontend
	if r.URL.Path == "/api/status" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "OK", "message": "WAF está operacional"})
		return
	}
	if r.URL.Path == "/api/logs" {
		mu.Lock()
		defer mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(blockedLogs)
		return
	}
	if r.URL.Path == "/" {
		http.ServeFile(w, r, "index.html")
		return
	}

	// Lógica do WAF (Proteção)
	if r.URL.Query().Get("user_input") == "drop+table" || r.URL.Path == "/evil" {
		log.Println("Requisição maliciosa detectada! Bloqueando...")
		
		// Registra o log
		mu.Lock()
		blockedLogs = append(blockedLogs, LogEntry{
			Timestamp:  time.Now().Format(time.RFC3339),
			RemoteAddr: r.RemoteAddr,
			Path:       r.URL.Path,
			Reason:     "Ataque de SQL Injection",
		})
		mu.Unlock()

		http.Error(w, "Requisição bloqueada por WAF!", http.StatusForbidden)
		return
	}

	// Proxy reverso para a aplicação real
	targetURL, err := url.Parse("http://localhost:8080")
	if err != nil {
		http.Error(w, "Erro interno do WAF", http.StatusInternalServerError)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)
}

func main() {
	http.HandleFunc("/", WAFHandler)
	log.Println("WAF iniciado e escutando na porta 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

