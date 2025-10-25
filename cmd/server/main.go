package main

import (
	"log"
	"net/http"
	"simple_zkp_integration/internal/handlers"

	"github.com/gorilla/mux"
)

func main() {

	// Initialize the router
	router := mux.NewRouter()

	// API routes
	router.HandleFunc("/api/proof/generate", handlers.GenerateProof).Methods("POST")
	router.HandleFunc("/api/proof/verify", handlers.VerifyProof).Methods("POST")
	router.HandleFunc("/health", handlers.HealthCheck).Methods("GET")

	router.Use(corsMiddleware)

	// Start server
	log.Println("🚀 ZKP API Server starting on port 8080...")
	log.Println("📝 API Endpoints:")
	log.Println("   POST /api/proof/generate - Generate a zero-knowledge proof")
	log.Println("   POST /api/proof/verify   - Verify a proof")
	log.Println("   GET  /health             - Health check")
	
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}