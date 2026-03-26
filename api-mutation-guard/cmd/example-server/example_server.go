// Example server: POST /payments, /deploy, /update require valid OER headers and body binding.
package main

import (
	"fmt"
	"log"
	"net/http"

	apimutationguard "github.com/TrigGuard-AI/TrigGuard/integrations/api-mutation-guard"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/payments", apimutationguard.Middleware(http.HandlerFunc(handlePayments)))
	mux.Handle("/deploy", apimutationguard.Middleware(http.HandlerFunc(handleDeploy)))
	mux.Handle("/update", apimutationguard.Middleware(http.HandlerFunc(handleUpdate)))

	addr := ":8080"
	log.Printf("example API mutation guard listening on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

func handlePayments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, `{"ok":true,"route":"payments"}`)
}

func handleDeploy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, `{"ok":true,"route":"deploy"}`)
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprint(w, `{"ok":true,"route":"update"}`)
}
