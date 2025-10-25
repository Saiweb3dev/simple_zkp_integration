// pkg/models/models.go
package models

// ProofRequest represents the request to generate a proof
type ProofRequest struct {
	A   int `json:"a"`   // First secret number
	B   int `json:"b"`   // Second secret number
	Sum int `json:"sum"` // Public sum (A + B)
}

// ProofResponse represents the response after generating a proof
type ProofResponse struct {
	Proof   []byte `json:"proof"`   // The zero-knowledge proof (serialized)
	Sum     int    `json:"sum"`     // The public sum
	Message string `json:"message"` // Human-readable message
}

// VerifyRequest represents the request to verify a proof
type VerifyRequest struct {
	Proof []byte `json:"proof"` // The proof to verify
	Sum   int    `json:"sum"`   // The claimed sum
}

// VerifyResponse represents the response after verifying a proof
type VerifyResponse struct {
	Valid   bool   `json:"valid"`   // Whether the proof is valid
	Message string `json:"message"` // Human-readable message
}