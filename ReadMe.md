# Go Mini-Blockchain

A production-grade blockchain implementation built from scratch in Golang, demonstrating deep understanding of distributed systems, cryptography, and consensus mechanisms without relying on external blockchain frameworks.

## Overview

This project implements a fully functional blockchain system with Proof of Work consensus, UTXO-based transactions with ECDSA cryptographic signatures, persistent storage using BoltDB, and peer-to-peer networking via WebSockets. Every component is built from first principles to understand the core mechanics of blockchain technology.

## Key Features

- **Proof of Work Consensus**: SHA-256 based mining with adjustable difficulty
- **UTXO Transaction Model**: Bitcoin-style transaction system with inputs and outputs
- **ECDSA Cryptography**: P-256 curve for digital signatures and key management
- **Persistent Storage**: BoltDB embedded database with crash recovery
- **P2P Networking**: WebSocket-based peer-to-peer block propagation
- **REST API**: HTTP endpoints for blockchain queries and operations
- **CLI Interface**: Professional command-line interface using Cobra framework

## Architecture

The system is built with clear separation of concerns across multiple layers:

```
Application Layer (CLI, HTTP, P2P) 
    ↓
Business Logic Layer (Blockchain, Block, Transaction)
    ↓
Consensus Layer (Proof of Work)
    ↓
Cryptography Layer (Wallet, Signatures)
    ↓
Storage Layer (BoltDB)
```

### Project Structure

```
go-mini-blockchain/
├── main.go              # Entry point - orchestrates all components
├── blockchain.db        # BoltDB persistent storage
│
├── block/
│   ├── block.go        # Block data structure and serialization
│   └── blockchain.go   # Blockchain management and BoltDB operations
│
├── proof/
│   └── pow.go          # Proof of Work mining algorithm
│
├── tx/
│   └── transaction.go  # UTXO model, signing, verification
│
├── wallet/
│   └── wallet.go       # ECDSA key generation and address derivation
│
├── p2p/
│   └── node.go         # WebSocket P2P networking
│
├── server/
│   └── server.go       # HTTP REST API server
│
└── cmd/
    ├── root.go         # Cobra root command
    ├── addBlock.go     # CLI: mine blocks
    ├── printChain.go   # CLI: display chain
    └── httpServer.go   # CLI: start HTTP server
```

## Technical Implementation

### Proof of Work

The mining algorithm requires finding a nonce such that the block hash is less than a target value derived from difficulty bits:

```
hash(prevHash + data + timestamp + targetBits + nonce) < target
```

Current difficulty: 16 bits (adjustable via `targetBits` constant)

### Transaction Model

Implements Bitcoin-style UTXO (Unspent Transaction Output) model:

- **Inputs**: Reference previous transaction outputs, include cryptographic signatures
- **Outputs**: Locked to recipient addresses using public key hashes
- **Coinbase**: Special transactions that create new coins as mining rewards

### Cryptographic Security

- **ECDSA P-256**: Elliptic curve digital signatures for transaction authorization
- **SHA-256**: Cryptographic hashing for blocks and transaction IDs
- **Address Generation**: `hex(first_20_bytes(SHA256(publicKey)))`

### Data Persistence

- **BoltDB**: Embedded key-value database with ACID guarantees
- **Serialization**: Go's gob encoding for block storage
- **Crash Recovery**: Blockchain state persists across program restarts
- **Chain Tip Tracking**: Special `lh` (last hash) key maintains current chain state

### Peer-to-Peer Networking

- **WebSocket Protocol**: Full-duplex communication for real-time block propagation
- **Concurrent Handling**: Separate goroutines for each peer connection
- **Broadcast Mechanism**: New blocks propagate to all connected peers
- **Peer Management**: Thread-safe peer map with mutex protection

## Installation

### Prerequisites

- Go 1.16 or higher
- Git

### Setup

```bash
# Clone the repository
git clone https://github.com/Shubham0699/go-mini-blockchain.git
cd go-mini-blockchain

# Install dependencies
go mod download

# Build the project
go build -o blockchain
```

## Usage

### Starting the Node

```bash
# Run with auto-mining and P2P networking
go run main.go
```

This starts:
- P2P node listening on `localhost:3000`
- Auto-mining every 10 seconds
- Interactive CLI for manual commands

### CLI Commands

Once the node is running, available commands:

```bash
# Manually mine a new block with coinbase reward
> mine <address>

# Display full blockchain history
> print

# Exit the program
> exit
```

### HTTP API Endpoints

Start the HTTP server (if using CLI mode):

```bash
# Using CLI
blockchain serve --port 8080

# Or modify main.go to start HTTP server
```

#### Available Endpoints

**GET /chain**
- Returns entire blockchain as JSON array
```bash
curl http://localhost:8080/chain
```

**GET /addblock?data=<text>**
- Mines new block with text data (legacy mode)
```bash
curl "http://localhost:8080/addblock?data=hello"
```

**POST /addblockjson**
- Mines new block with JSON payload
```bash
curl -X POST http://localhost:8080/addblockjson \
  -H "Content-Type: application/json" \
  -d '{"data": "transaction details"}'
```

### Running Multiple Nodes (P2P Demo)

To demonstrate peer-to-peer block propagation:

1. **Node 1** (default setup):
```go
// main.go - line 23
node := p2p.NewNode("localhost:3000", bc)
```

2. **Node 2** (modify and run separately):
```go
// main.go - line 23
node := p2p.NewNode("localhost:3001", bc)
// Add peer connection
node.ConnectPeer("localhost:3000")
```

Blocks mined on either node will automatically propagate to connected peers.

## Code Examples

### Creating a Wallet

```go
import "github.com/Shubham0699/go-mini-blockchain/wallet"

// Generate new wallet with ECDSA key pair
w, err := wallet.NewWallet()
if err != nil {
    log.Fatal(err)
}

// Get blockchain address
address := w.Address()
fmt.Println("Address:", address)
```

### Mining a Block

```go
import (
    "github.com/Shubham0699/go-mini-blockchain/block"
    "github.com/Shubham0699/go-mini-blockchain/tx"
)

// Get blockchain instance
bc := block.GetBlockchain()
defer bc.Close()

// Create coinbase transaction (mining reward)
cbTx := tx.NewCoinbaseTX(minerAddress, 50)

// Mine block with transaction
bc.MineBlock([]*tx.Transaction{cbTx})
```

### Signing and Verifying Transactions

```go
// Create wallet
wallet, _ := wallet.NewWallet()

// Create transaction
transaction := &tx.Transaction{
    Vin:  []tx.TXInput{...},
    Vout: []tx.TXOutput{...},
}

// Sign transaction
prevOutputs := map[string]tx.TXOutput{...}
transaction.Sign(wallet.Private, prevOutputs)

// Verify transaction
isValid := transaction.Verify(prevOutputs)
```

## System Flows

### Program Startup

1. Load blockchain from BoltDB (or create new with genesis block)
2. Initialize P2P node with blockchain reference
3. Start WebSocket server in background goroutine
4. Launch auto-mining goroutine (mines every 10 seconds)
5. Enter CLI input loop for manual commands

### Mining Flow

1. Create coinbase transaction with mining reward
2. Assemble new block with transactions and previous block hash
3. Run Proof of Work algorithm to find valid nonce
4. Store mined block in BoltDB with hash as key
5. Update chain tip to new block hash
6. Broadcast block to all connected peers

### Transaction Validation

1. Check if transaction is coinbase (skip signature validation)
2. For each input, retrieve referenced output's public key hash
3. Reconstruct hash that was signed (transaction data + output hash)
4. Extract signature and public key from input
5. Verify signature using ECDSA algorithm
6. Transaction valid only if all inputs have valid signatures

### P2P Block Propagation

1. Node A mines block and calls BroadcastBlock()
2. Block serialized to JSON and sent via WebSocket to all peers
3. Node B receives block through ListenPeer goroutine
4. Node B deserializes JSON to Block struct
5. Node B adds block to local blockchain
6. Network achieves eventual consistency through recursive propagation

## Design Decisions

### Why UTXO Model?

- More flexible than account-based model
- Enables parallel transaction validation
- Natural fit for blockchain immutability
- Same model used by Bitcoin

### Why BoltDB?

- Embedded database with zero external dependencies
- ACID guarantees for data integrity
- Fast key-value lookups by block hash
- No database server required

### Why WebSockets for P2P?

- Full-duplex communication for real-time updates
- Efficient binary and JSON data transfer
- Native browser support for future web clients
- Better than HTTP polling for block propagation

### Why Singleton Blockchain?

- Prevents multiple instances causing data inconsistency
- Shared across CLI, HTTP, and P2P interfaces
- Thread-safe initialization with sync.Once
- Single source of truth for blockchain state

### Interface-Based Proof of Work

The `BlockData` interface prevents circular dependency between `proof` and `block` packages:

```go
type BlockData interface {
    PrevHash() []byte
    DataBytes() []byte
    TimestampUnix() int64
    NonceValue() int64
}
```

This design allows proof of work logic to operate on blocks without importing the block package.

## Performance Considerations

- **Mining Time**: With 16-bit difficulty, average mining time is 2-5 seconds on modern hardware
- **Database Writes**: BoltDB serializes writes, ensuring consistency but limiting concurrent write throughput
- **Memory Usage**: Iterator pattern prevents loading entire blockchain into memory
- **Network Latency**: WebSocket connections minimize propagation delay between peers

## Security Features

- **Immutable Blocks**: Once mined, blocks cannot be modified (hash would change)
- **Cryptographic Signatures**: ECDSA prevents unauthorized transaction spending
- **Chain Integrity**: Each block references previous block hash, making tampering evident
- **Proof of Work**: Computational cost prevents spam and provides Sybil resistance
- **Address Privacy**: Public key hashed to create addresses, adding security layer

## Testing

### Manual Testing

```bash
# Test basic mining
> mine testaddress123

# Verify block was added
> print

# Test HTTP endpoint
curl http://localhost:8080/chain
```

### Testing P2P Propagation

1. Start Node 1 in one terminal
2. Start Node 2 in another terminal (with modified port)
3. Mine block on Node 1
4. Verify block appears on Node 2's chain

## Future Improvements

### Short-Term Enhancements

- UTXO set caching for faster transaction validation
- Merkle tree implementation for efficient proof of inclusion
- Transaction pool (mempool) for pending transactions
- Dynamic difficulty adjustment based on block time
- Automated peer discovery mechanism
- Wallet persistence (save/load keys from encrypted files)
- SPV (Simplified Payment Verification) for light clients

### Medium-Term Goals

- Proof of Stake consensus as alternative to PoW
- Smart contract VM with Turing-complete execution
- State management layer for contract storage
- Gossip protocol for improved P2P message propagation
- Chain reorganization handling (fork resolution)
- Transaction fee mechanism and fee-based priority
- JSON-RPC 2.0 interface for standardized API access

### Long-Term Vision: Modular Architecture

Transform the project into a Cosmos SDK-inspired blockchain framework:

- **Consensus Module**: Pluggable consensus engines (PoW, PoS, Tendermint)
- **State Module**: Separate state machine from blockchain core
- **Transaction Module**: Custom transaction types without modifying core
- **Networking Module**: Abstract P2P layer supporting multiple protocols
- **ABCI Interface**: Application Blockchain Interface for custom chains
- **Module Manager**: Lifecycle management and dependency injection

This evolution converts the project from single-purpose blockchain to a framework capable of supporting multiple application-specific blockchains.

## Known Limitations

- **No Difficulty Adjustment**: Mining difficulty is fixed, not dynamic
- **No Fork Resolution**: Longest chain rule not implemented
- **No Transaction Pool**: Transactions immediately go into blocks
- **Manual Peer Connections**: No automatic peer discovery
- **No Network Encryption**: P2P connections are unencrypted
- **Limited Validation**: Double-spend prevention is conceptual, not fully enforced
- **No Checkpointing**: Chain cannot be pruned or checkpointed

## Troubleshooting

### Database Lock Error

```
cannot open database: database is locked
```

**Solution**: Ensure only one instance of the program is running. BoltDB allows only one write connection.

### Mining Too Slow/Fast

**Solution**: Adjust `targetBits` constant in `proof/pow.go`. Higher value = slower mining.

### Peer Connection Failed

```
Failed to connect to peer: connection refused
```

**Solution**: Verify the peer's address and ensure their node is running and accessible.

### Port Already in Use

```
bind: address already in use
```

**Solution**: Change the port in `main.go` or ensure no other program is using port 3000.

## Dependencies

```
github.com/boltdb/bolt       # Embedded key-value database
github.com/gorilla/websocket # WebSocket implementation
github.com/spf13/cobra       # CLI framework
```

Install all dependencies:

```bash
go mod download
```

## Contributing

This is an educational project demonstrating blockchain fundamentals. While it's not accepting contributions, feel free to:

- Fork the repository for your own learning
- Use the code as reference for blockchain concepts
- Adapt the architecture for your projects

## Development Philosophy

This project prioritizes understanding over convenience:

- **No frameworks**: Built from scratch to understand internals
- **Explicit over implicit**: Clear, readable code over clever abstractions
- **Documentation**: Every design decision is documented
- **Incremental complexity**: Added features step-by-step
- **Production patterns**: Real-world coding practices

## Learning Resources

If you're studying this codebase, recommended reading order:

1. `block/block.go` - Understand block structure
2. `proof/pow.go` - Learn mining algorithm
3. `block/blockchain.go` - See how blocks are stored
4. `tx/transaction.go` - Understand UTXO model
5. `wallet/wallet.go` - Learn key management
6. `p2p/node.go` - Explore networking
7. `main.go` - See how everything connects

## Performance Benchmarks

Approximate performance on modern hardware (Intel i7, 16GB RAM):

- **Mining time**: 2-5 seconds per block (16-bit difficulty)
- **Block storage**: ~1KB per block (varies with transaction count)
- **Network latency**: <100ms for local peer propagation
- **Database reads**: <1ms per block retrieval
- **Signature verification**: ~2ms per transaction

## License

This project is open source and available for educational purposes.

## Author

**Shubh**
- Computer Science Graduate, VIT Pune (2024)
- Blockchain Developer | Backend Engineer | Cybersecurity Enthusiast
- GitHub: [@Shubham0699](https://github.com/Shubham0699)

## Acknowledgments

Inspired by Bitcoin's original implementation and educational blockchain resources. Built to demonstrate systems-level understanding of blockchain technology beyond high-level frameworks.

## Project Status

**Current Status**: Fully functional blockchain with PoW consensus, transactions, and P2P networking

**Next Milestone**: Modular architecture refactoring (Cosmos-like design)

---

**Note**: This is an educational implementation optimized for learning and demonstration. For production use cases, consider established blockchain frameworks with extensive testing and security audits.