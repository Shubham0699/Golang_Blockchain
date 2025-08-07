package block

import (
    "bytes"
    "fmt"

    "github.com/Shubham0699/go-mini-blockchain/proof"
)

type Blockchain struct {
    Blocks []*Block // ✅ Capitalized to export
}

// AddBlock adds a new block to the chain
func (bc *Blockchain) AddBlock(data string) {
    prevBlock := bc.Blocks[len(bc.Blocks)-1]
    newBlock := NewBlock(data, prevBlock.Hash)
    bc.Blocks = append(bc.Blocks, newBlock)
}

// NewGenesisBlock creates the first block
func NewGenesisBlock() *Block {
    return NewBlock("genesis block", []byte{})
}

// NewBlockchain creates a new blockchain with the genesis block
func NewBlockchain() *Blockchain {
    return &Blockchain{[]*Block{NewGenesisBlock()}}
}

// ✅ IsValid checks if the blockchain is valid
func (bc *Blockchain) IsValid() bool {
    for i := 1; i < len(bc.Blocks); i++ {
        current := bc.Blocks[i]
        previous := bc.Blocks[i-1]

        pow := proof.NewProofOfWork(current)
        if !pow.Validate() {
            fmt.Printf("❌ Invalid Proof of Work at block %d\n", i)
            return false
        }

        if !bytes.Equal(current.PrevBlockHash, previous.Hash) {
            fmt.Printf("❌ Invalid PrevHash linkage at block %d\n", i)
            return false
        }
    }

    fmt.Println("✅ Blockchain is valid.")
    return true
}
