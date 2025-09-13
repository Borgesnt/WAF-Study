package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

// WAFHandler é o nosso handler principal que atua como proxy
func WAFHandler(w http.ResponseWriter, r *http.Request) {
	// Passo 1: Recebe a requisição
	fmt.Printf("Recebendo requisição de: %s\n", r.RemoteAddr)

	// Passo 2: Analisa a requisição (lógica do WAF)
	// **Aqui você implementaria a lógica para inspecionar a requisição
	// por padrões maliciosos, como SQL Injection ou XSS.**
	// Por enquanto, vamos apenas logar.
	if r.URL.Path == "/evil" {
		log.Println("Requisição maliciosa detectada! Bloqueando...")
		// Passo 3: Bloqueia a requisição e retorna um erro
		http.Error(w, "Requisição bloqueada por WAF!", http.StatusForbidden)
		return
	}

	// Se a requisição não for maliciosa, encaminha para a aplicação
	// Passo 4: Encaminha a requisição
	targetURL, _ := url.Parse("http://localhost:8080") // Alvo é a sua aplicação real
	proxy := httputil.NewSingleHostReverseProxy(targetURL)
	proxy.ServeHTTP(w, r)
}

func main() {
	// O WAF vai escutar na porta 8000 do container
	http.HandleFunc("/", WAFHandler)
	log.Println("WAF iniciado e escutando na porta 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
