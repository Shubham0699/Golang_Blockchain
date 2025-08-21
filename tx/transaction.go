package tx

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"math/big"
)

// TXInput represents a transaction input
type TXInput struct {
	Txid      []byte
	Vout      int
	Signature []byte
	PubKey    []byte
}

// TXOutput represents a transaction output
type TXOutput struct {
	Value      int
	PubKeyHash []byte // locked to an address-hash
}

func (out *TXOutput) Lock(address string) {
	// address is hex(20 bytes). Decode to raw 20 bytes.
	b, _ := hex.DecodeString(address)
	out.PubKeyHash = b
}

func NewTXOutput(value int, address string) TXOutput {
	o := TXOutput{Value: value}
	o.Lock(address)
	return o
}

// Transaction holds inputs and outputs
type Transaction struct {
	ID   []byte
	Vin  []TXInput
	Vout []TXOutput
}

func (tx *Transaction) Hash() []byte {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	_ = enc.Encode(tx.Vin)
	_ = enc.Encode(tx.Vout)
	h := sha256.Sum256(buf.Bytes())
	return h[:]
}

func (tx *Transaction) SetID() { tx.ID = tx.Hash() }

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// NewCoinbaseTX creates a coinbase (mining reward) transaction
func NewCoinbaseTX(to string, reward int) *Transaction {
	tx := &Transaction{
		Vin:  []TXInput{{Txid: []byte{}, Vout: -1, Signature: nil, PubKey: nil}},
		Vout: []TXOutput{NewTXOutput(reward, to)},
	}
	tx.SetID()
	return tx
}

// Sign signs each input of the transaction with the provided private key.
// prevOutMap maps "txid||vout" (hex encoded) to the referenced TXOutput.
func (tx *Transaction) Sign(priv *ecdsa.PrivateKey, prevOutMap map[string]TXOutput) {
	if tx.IsCoinbase() {
		return
	}
	txCopy := tx.trimmedCopy()

	for inIdx := range txCopy.Vin {
		in := &txCopy.Vin[inIdx]
		// include the referenced output's PubKeyHash into the hash
		key := hex.EncodeToString(append(in.Txid, byte(in.Vout)))
		prevOut := prevOutMap[key]

		// hash: txCopy + referenced output PKH
		h := sha256.Sum256(append(txCopy.Hash(), prevOut.PubKeyHash...))

		r, s, err := ecdsa.Sign(rand.Reader, priv, h[:])
		if err != nil {
			// In production you'd return an error; for now panic to keep behavior consistent with the rest of the project
			panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)

		// write signature + pubkey back to original tx
		tx.Vin[inIdx].Signature = signature
		tx.Vin[inIdx].PubKey = append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	}
}

// Verify verifies signatures of transaction inputs using prevOutMap (same format as Sign)
func (tx *Transaction) Verify(prevOutMap map[string]TXOutput) bool {
	if tx.IsCoinbase() {
		return true
	}
	txCopy := tx.trimmedCopy()

	for inIdx, vin := range tx.Vin {
		key := hex.EncodeToString(append(vin.Txid, byte(vin.Vout)))
		prevOut := prevOutMap[key]

		h := sha256.Sum256(append(txCopy.Hash(), prevOut.PubKeyHash...))

		// restore r,s
		sig := vin.Signature
		if len(sig) == 0 {
			return false
		}
		r := new(big.Int).SetBytes(sig[:len(sig)/2])
		s := new(big.Int).SetBytes(sig[len(sig)/2:])

		// restore pubkey
		px := new(big.Int).SetBytes(vin.PubKey[:len(vin.PubKey)/2])
		py := new(big.Int).SetBytes(vin.PubKey[len(vin.PubKey)/2:])
		pub := ecdsa.PublicKey{Curve: elliptic.P256(), X: px, Y: py}

		if !ecdsa.Verify(&pub, h[:], r, s) {
			return false
		}

		// wipe this field in the copy like standard signing scheme (not strictly required here)
		txCopy.Vin[inIdx].Signature = nil
		txCopy.Vin[inIdx].PubKey = nil
	}
	return true
}

func (tx *Transaction) trimmedCopy() *Transaction {
	var inputs []TXInput
	for _, in := range tx.Vin {
		inputs = append(inputs, TXInput{Txid: in.Txid, Vout: in.Vout})
	}
	var outputs []TXOutput
	outputs = append(outputs, tx.Vout...)
	return &Transaction{Vin: inputs, Vout: outputs}
}
