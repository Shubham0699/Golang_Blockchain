package block

import (
    "bytes"
    "crypto/sha256"
    "encoding/gob"
    "log"
    "strconv"
    "time"

    "github.com/Shubham0699/go-mini-blockchain/proof"
    "github.com/Shubham0699/go-mini-blockchain/tx"
)

type Block struct {
    Timestamp     int64
    Data          []byte
    PrevBlockHash []byte
    Hash          []byte
    Nonce         int64
    Transactions  []*tx.Transaction
}

func init() {
    gob.Register(&tx.Transaction{})
}

// Implementing proof.BlockData interface
func (b *Block) PrevHash() []byte        { return b.PrevBlockHash }
func (b *Block) DataBytes() []byte       { return b.Data }
func (b *Block) TimestampUnix() int64    { return b.Timestamp }
func (b *Block) NonceValue() int64       { return b.Nonce }

// Legacy SetHash (not used with PoW)
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

// ✅ NewBlock for simple string data blocks
func NewBlock(data string, prevBlockHash []byte) *Block {
    block := &Block{
        Timestamp:     time.Now().Unix(),
        Data:          []byte(data),
        PrevBlockHash: prevBlockHash,
        Hash:          []byte{},
        Nonce:         0,
        Transactions:  nil,
    }

    pow := proof.NewProofOfWork(block)
    nonce, hash := pow.Run()
    block.Hash = hash
    block.Nonce = nonce

    return block
}

// NewBlockWithTxs creates a block containing transactions
func NewBlockWithTxs(transactions []*tx.Transaction, prevBlockHash []byte) *Block {
    block := &Block{
        Timestamp:     time.Now().Unix(),
        Data:          nil,
        PrevBlockHash: prevBlockHash,
        Hash:          []byte{},
        Nonce:         0,
        Transactions:  transactions,
    }

    pow := proof.NewProofOfWork(block)
    nonce, hash := pow.Run()
    block.Hash = hash
    block.Nonce = nonce

    return block
}

// Genesis block
func NewGenesisBlock() *Block {
    return NewBlock("Genesis Block", []byte{})
}

// Serialize block to bytes
func (b *Block) Serialize() []byte {
    var result bytes.Buffer
    encoder := gob.NewEncoder(&result)

    err := encoder.Encode(b)
    if err != nil {
        log.Panic(err)
    }

    return result.Bytes()
}

// Deserialize bytes to block
func Deserialize(d []byte) *Block {
    var block Block

    decoder := gob.NewDecoder(bytes.NewReader(d))
    err := decoder.Decode(&block)
    if err != nil {
        log.Panic(err)
    }

    return &block
}
