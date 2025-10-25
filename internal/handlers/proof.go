// internal/handlers/proof.go
package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"simple_zkp_integration/internal/circuit"
	"simple_zkp_integration/pkg/models"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark-crypto/ecc"
)

var (
	// Global keys stored in memory (in production, these would be persisted)
	provingKey   groth16.ProvingKey
	verifyingKey groth16.VerifyingKey
	setupOnce    sync.Once
	setupErr     error
)

// initializeKeys performs the one-time setup to generate keys
func initializeKeys() {
	setupOnce.Do(func() {
		log.Println("🔧 Performing trusted setup (generating keys)...")
		provingKey, verifyingKey, setupErr = circuit.Setup()
		if setupErr != nil {
			log.Printf("❌ Setup failed: %v", setupErr)
			return
		}
		log.Println("✅ Setup complete! Ready to generate and verify proofs.")
	})
}

// GenerateProof handles POST /api/proof/generate
func GenerateProof(w http.ResponseWriter, r *http.Request) {
	// Initialize keys if not already done
	initializeKeys()
	if setupErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Setup failed: "+setupErr.Error())
		return
	}

	// Parse request
	var req models.ProofRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Optional validation - the circuit will enforce this constraint anyway
	// You can uncomment this for early validation if desired
	// if req.A+req.B != req.Sum {
	// 	respondWithError(w, http.StatusBadRequest, "Invalid inputs: A + B must equal Sum")
	// 	return
	// }

	log.Printf("📝 Generating proof for: %d + %d = %d", req.A, req.B, req.Sum)

	// Generate proof
	proof, err := circuit.GenerateProof(provingKey, req.A, req.B, req.Sum)
	if err != nil {
		log.Printf("❌ Proof generation failed: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to generate proof")
		return
	}

	// Serialize proof to bytes using WriteTo
	// This is the standard way to serialize proofs in gnark
	buf := new(bytes.Buffer)
	if _, err := proof.WriteTo(buf); err != nil {
		log.Printf("❌ Proof serialization failed: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Failed to serialize proof")
		return
	}
	proofBytes := buf.Bytes()

	log.Printf("✅ Proof generated successfully!")

	// Return response
	response := models.ProofResponse{
		Proof:   proofBytes,
		Sum:     req.Sum,
		Message: "Proof generated successfully. This proves you know two numbers that add up to the sum, without revealing the numbers themselves!",
	}

	respondWithJSON(w, http.StatusOK, response)
}

// VerifyProof handles POST /api/proof/verify
func VerifyProof(w http.ResponseWriter, r *http.Request) {
	// Initialize keys if not already done
	initializeKeys()
	if setupErr != nil {
		respondWithError(w, http.StatusInternalServerError, "Setup failed: "+setupErr.Error())
		return
	}

	// Parse request
	var req models.VerifyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	log.Printf("🔍 Verifying proof for sum: %d", req.Sum)

	// Deserialize proof
	proof := groth16.NewProof(ecc.BN254)
	buf := bytes.NewReader(req.Proof)
	if _, err := proof.ReadFrom(buf); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid proof format")
		return
	}

	// Verify proof
	err := circuit.VerifyProof(verifyingKey, proof, req.Sum)
	
	if err != nil {
		log.Printf("❌ Proof verification failed: %v", err)
		response := models.VerifyResponse{
			Valid:   false,
			Message: "Proof is invalid. The prover does not know valid numbers that add up to the given sum.",
		}
		respondWithJSON(w, http.StatusOK, response)
		return
	}

	log.Printf("✅ Proof verified successfully!")

	response := models.VerifyResponse{
		Valid:   true,
		Message: "Proof is valid! The prover knows two numbers that add up to the sum (without revealing them).",
	}

	respondWithJSON(w, http.StatusOK, response)
}

// HealthCheck handles GET /health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{
		"status":  "healthy",
		"service": "zkp-api",
	})
}

// Helper functions for JSON responses
func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}