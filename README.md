# Zero-Knowledge Proof API Tutorial

A beginner-friendly guide to integrating Zero-Knowledge Proofs into your Go backend using gnark.

## 🎯 What This Project Does

This API demonstrates a simple zero-knowledge proof system where a user can prove they know two numbers that add up to a specific sum **without revealing the actual numbers**. This is the fundamental concept behind ZKP technology used in blockchain privacy solutions, authentication systems, and secure computation.

## 📁 Project Structure

```
zkp-api-tutorial/
├── cmd/
│   └── server/
│       └── main.go              # API server entry point
├── internal/
│   ├── circuit/
│   │   └── addition.go          # ZKP circuit implementation
│   └── handlers/
│       └── proof.go             # HTTP request handlers
├── pkg/
│   └── models/
│       └── models.go            # Request/response models
├── Dockerfile                   # Container definition
├── docker-compose.yml           # Container orchestration
├── go.mod                       # Go dependencies
└── README.md                    # This file
```

## 🚀 Quick Start (5 minutes)

### Prerequisites

- Docker and Docker Compose installed
- OR: Go 1.21+ installed (for local development)

### Option 1: Using Docker (Recommended)

```bash
# Clone the repository
git clone https://github.com/yourusername/zkp-api-tutorial.git
cd zkp-api-tutorial

# Start the server
docker-compose up --build

# The API will be available at http://localhost:8080
```

### Option 2: Local Development

```bash
# Clone the repository
git clone https://github.com/yourusername/zkp-api-tutorial.git
cd zkp-api-tutorial

# Install dependencies
go mod download

# Run the server
go run cmd/server/main.go

# The API will be available at http://localhost:8080
```

## 📡 API Endpoints

### 1. Health Check
```bash
curl http://localhost:8080/health
```

### 2. Generate Proof
Prove you know two numbers that add up to a sum without revealing them:

```bash
curl -X POST http://localhost:8080/api/proof/generate \
  -H "Content-Type: application/json" \
  -d '{
    "a": 5,
    "b": 3,
    "sum": 8
  }'
```

**Response:**
```json
{
  "proof": "0x1a2b3c...",
  "sum": 8,
  "message": "Proof generated successfully..."
}
```

### 3. Verify Proof
Verify a proof without knowing the secret numbers:

```bash
curl -X POST http://localhost:8080/api/proof/verify \
  -H "Content-Type: application/json" \
  -d '{
    "proof": "0x1a2b3c...",
    "sum": 8
  }'
```

**Response:**
```json
{
  "valid": true,
  "message": "Proof is valid! The prover knows two numbers that add up to the sum."
}
```

## 🧪 Testing the API

Here's a complete workflow to test the ZKP system:

```bash
# Step 1: Generate a proof (you know a=5 and b=3)
RESPONSE=$(curl -s -X POST http://localhost:8080/api/proof/generate \
  -H "Content-Type: application/json" \
  -d '{"a": 5, "b": 3, "sum": 8}')

echo "Generated proof: $RESPONSE"

# Step 2: Extract the proof from response
PROOF=$(echo $RESPONSE | jq -r '.proof')

# Step 3: Verify the proof (only the sum is public)
curl -X POST http://localhost:8080/api/proof/verify \
  -H "Content-Type: application/json" \
  -d "{\"proof\": $PROOF, \"sum\": 8}"
```

## 🔧 Development

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
# Build the binary
go build -o server ./cmd/server

# Run the binary
./server
```

### Docker Build
```bash
# Build the image
docker build -t zkp-api .

# Run the container
docker run -p 8080:8080 zkp-api
```

## 📚 Understanding the Code

### The Circuit (`internal/circuit/addition.go`)
The circuit defines the mathematical relationship we want to prove. In this case: `A + B = Sum`

The circuit has:
- **Private inputs** (A, B): Secret values known only to the prover
- **Public input** (Sum): Value known to everyone
- **Constraint**: The assertion that A + B must equal Sum

### The API Flow
1. **Setup Phase** (happens once on server start):
   - Generates proving and verifying keys
   - This is like creating a "lock" and "key" for the cryptographic system

2. **Proof Generation** (`/api/proof/generate`):
   - Takes secret values (a, b) and public sum
   - Creates cryptographic proof that a + b = sum
   - Returns proof without revealing a and b

3. **Proof Verification** (`/api/proof/verify`):
   - Takes proof and public sum
   - Verifies proof is valid
   - Confirms prover knows correct values WITHOUT seeing them

## 🎓 Learning Resources

- [gnark Documentation](https://docs.gnark.consensys.net/)
- [Zero-Knowledge Proofs Explained](https://consensys.net/blog/blockchain-explained/zero-knowledge-proofs-starks-vs-snarks/)
- [ZK-SNARKs Tutorial](https://github.com/matter-labs/awesome-zero-knowledge-proofs)

## 🤝 Contributing

Feel free to open issues or submit pull requests to improve this tutorial!

## 📄 License

MIT License - feel free to use this for learning and teaching.

## ⚠️ Production Considerations

This is a **tutorial project** for learning. Before using ZKP in production:

1. **Trusted Setup**: The setup phase should use a secure multi-party computation ceremony
2. **Key Management**: Store proving/verifying keys securely (not in memory)
3. **Circuit Auditing**: Have cryptography experts audit your circuits
4. **Performance**: Consider circuit optimization for complex operations
5. **Error Handling**: Add comprehensive error handling and logging

## 💡 Next Steps

Once you understand this basic example, try:
- Creating more complex circuits (multiplication, comparison)
- Adding authentication with ZKP
- Implementing range proofs
- Building a blockchain integration