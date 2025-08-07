package block

import (
    "bytes"
    "crypto/sha256"
    "strconv"
    "time"

    "github.com/Shubham0699/go-mini-blockchain/proof"
)

type Block struct {
    Timestamp     int64
    Data          []byte
    PrevBlockHash []byte
    Hash          []byte
    Nonce         int64
}

// Implementing proof.BlockData interface:

func (b *Block) PrevHash() []byte {
    return b.PrevBlockHash
}

func (b *Block) DataBytes() []byte {
    return b.Data
}

func (b *Block) TimestampUnix() int64 {
    return b.Timestamp
}

func (b *Block) NonceValue() int64 {
    return b.Nonce
}

// SetHash is now unused because we use PoW, but you can keep it for reference
func (b *Block) SetHash() {
    headers := bytes.Join(
        [][]byte{
            b.PrevBlockHash,
            b.Data,
            []byte(strconv.FormatInt(b.Timestamp, 10)),
        },
        []byte{},
    )
    hash := sha256.Sum256(headers)
    b.Hash = hash[:]
}

// NewBlock creates a new block and mines it using PoW
func NewBlock(data string, prevBlockHash []byte) *Block {
    block := &Block{
        Timestamp:     time.Now().Unix(),
        Data:          []byte(data),
        PrevBlockHash: prevBlockHash,
        Hash:          []byte{},
        Nonce:         0,
    }

    pow := proof.NewProofOfWork(block) // block now implements BlockData
    nonce, hash := pow.Run()
    block.Hash = hash
    block.Nonce = nonce

    return block
}
