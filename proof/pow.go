package proof

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

// Target bits define the difficulty. More bits = harder.
const targetBits = 16

// 👇 This interface removes the need to import the block package
type BlockData interface {
	PrevHash() []byte
	DataBytes() []byte
	TimestampUnix() int64
	NonceValue() int64
}

type ProofOfWork struct {
	Block  BlockData
	Target *big.Int
}

// Prepare data for hashing
func (pow *ProofOfWork) prepareData(nonce int64) []byte {
	return bytes.Join(
		[][]byte{
			pow.Block.PrevHash(),
			pow.Block.DataBytes(),
			[]byte(fmt.Sprintf("%d", pow.Block.TimestampUnix())),
			[]byte(fmt.Sprintf("%d", targetBits)),
			[]byte(fmt.Sprintf("%d", nonce)),
		},
		[]byte{},
	)
}

// Main mining loop
func (pow *ProofOfWork) Run() (int64, []byte) {
	var hashInt big.Int
	var hash [32]byte
	var nonce int64 = 0

	fmt.Println("⛏️ Mining a new block...")

	for nonce < math.MaxInt64 {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.Target) == -1 {
			break // found!
		} else {
			nonce++
		}
	}

	fmt.Printf("✅ Mined! Nonce: %d\n", nonce)
	fmt.Printf("🔑 Hash: %x\n", hash[:])

	return nonce, hash[:]
}

// Validates PoW
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int
	data := pow.prepareData(pow.Block.NonceValue())
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	return hashInt.Cmp(pow.Target) == -1
}

// Constructor
func NewProofOfWork(b BlockData) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))
	return &ProofOfWork{b, target}
}
