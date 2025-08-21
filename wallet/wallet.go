package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

type Wallet struct {
	Private *ecdsa.PrivateKey
	PubKey  []byte // uncompressed: X||Y
}

func NewWallet() (*Wallet, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}
	pub := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	return &Wallet{Private: priv, PubKey: pub}, nil
}

// Very simple address: hex( first 20 bytes of SHA256(pubkey) )
// (We can upgrade to RIPEMD160+Base58Check later without touching call sites.)
func (w *Wallet) Address() string {
	h := sha256.Sum256(w.PubKey)
	return hex.EncodeToString(h[:20])
}
