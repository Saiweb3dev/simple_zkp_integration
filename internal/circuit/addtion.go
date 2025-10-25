package circuit

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark-crypto/ecc"
)

// AdditionCircuit defines a simple addition circuit
// This circuit proves: "I know two numbers (a, b) that add up to a given sum"
// without revealing what a and b actually are
type AdditionCircuit struct {
	// Private inputs (witness) - these are secret
	A frontend.Variable `gnark:",secret"`
	B frontend.Variable `gnark:",secret"`
	
	// Public input - this is known to everyone
	Sum frontend.Variable `gnark:",public"`
}

// Define declares the circuit's constraints
// This is where we define the mathematical relationship we want to prove
func (circuit *AdditionCircuit) Define(api frontend.API) error {
	// Assert that A + B = Sum
	// This creates a constraint that must be satisfied for the proof to be valid
	api.AssertIsEqual(api.Add(circuit.A, circuit.B), circuit.Sum)
	return nil
}

// Setup generates the proving and verifying keys
// This is a one-time setup phase that's needed before generating proofs
func Setup() (groth16.ProvingKey, groth16.VerifyingKey, error) {
	// Create an instance of our circuit
	var circuit AdditionCircuit
	
	// Compile the circuit into a constraint system
	// This converts our high-level circuit into mathematical constraints
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return nil, nil, err
	}
	
	// Run the trusted setup to generate keys
	// In production, this would be done through a secure ceremony
	pk, vk, err := groth16.Setup(ccs)
	if err != nil {
		return nil, nil, err
	}
	
	return pk, vk, nil
}

// GenerateProof creates a zero-knowledge proof for given inputs
func GenerateProof(pk groth16.ProvingKey, a, b, sum int) (groth16.Proof, error) {
	// Create witness (the actual values we want to prove)
	assignment := AdditionCircuit{
		A:   a,
		B:   b,
		Sum: sum,
	}

	// Create witness vector
	witness, err := frontend.NewWitness(&assignment, ecc.BN254.ScalarField())
	if err != nil {
		return nil, err
	}

	// Compile the circuit into a constraint system (needed by Prove)
	var circuit AdditionCircuit
	ccs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		return nil, err
	}

	// Generate the proof
	// This creates cryptographic evidence that a + b = sum
	// without revealing what a and b are
	proof, err := groth16.Prove(ccs, pk, witness)
	if err != nil {
		return nil, err
	}

	return proof, nil
}

// VerifyProof checks if a proof is valid
func VerifyProof(vk groth16.VerifyingKey, proof groth16.Proof, sum int) error {
	// Create public witness (only the sum is public)
	publicAssignment := AdditionCircuit{
		Sum: sum,
	}
	
	// Create public witness vector
	publicWitness, err := frontend.NewWitness(&publicAssignment, ecc.BN254.ScalarField(), frontend.PublicOnly())
	if err != nil {
		return err
	}
	
	// Verify the proof
	// This checks the cryptographic proof without knowing a and b
	err = groth16.Verify(proof, vk, publicWitness)
	return err
}