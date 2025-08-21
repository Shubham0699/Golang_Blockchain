package block

import (
	"log"
	"sync"

	"github.com/boltdb/bolt"
	"github.com/Shubham0699/go-mini-blockchain/tx"
)

const (
	dbFile       = "blockchain.db"
	blocksBucket = "blocks"
	lastHashKey  = "lh"
)

// Blockchain represents the chain stored in BoltDB
type Blockchain struct {
	tip []byte   // last block hash
	db  *bolt.DB // BoltDB instance
}

// CreateBlockchain creates a new blockchain with a genesis block
func CreateBlockchain() *Blockchain {
	var tip []byte

	db, err := bolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			// No existing chain → create one
			genesis := NewGenesisBlock()

			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			// Serialize genesis block
			err = b.Put(genesis.Hash, genesis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			// Save last hash
			err = b.Put([]byte(lastHashKey), genesis.Hash)
			if err != nil {
				log.Panic(err)
			}

			tip = genesis.Hash
		} else {
			// Chain exists → load last hash
			tip = b.Get([]byte(lastHashKey))
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	return &Blockchain{tip, db}
}

// AddBlock saves a new block into BoltDB (string data payload)
func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lastHash = b.Get([]byte(lastHashKey))
		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lastHash)

	err = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte(lastHashKey), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}

		bc.tip = newBlock.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// MineBlock mines a new block containing real transactions
func (bc *Blockchain) MineBlock(transactions []*tx.Transaction) {
	newBlock := NewBlockWithTxs(transactions, bc.tip)

	err := bc.db.Update(func(txn *bolt.Tx) error {
		b := txn.Bucket([]byte(blocksBucket))

		if err := b.Put(newBlock.Hash, newBlock.Serialize()); err != nil {
			log.Panic(err)
		}
		if err := b.Put([]byte(lastHashKey), newBlock.Hash); err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash
		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

// Iterator to traverse blockchain
type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	return &BlockchainIterator{bc.tip, bc.db}
}

func (it *BlockchainIterator) Next() *Block {
	var block *Block

	err := it.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(it.currentHash)
		block = Deserialize(encodedBlock)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	it.currentHash = block.PrevBlockHash
	return block
}

// Close closes the underlying BoltDB
func (bc *Blockchain) Close() {
	bc.db.Close()
}

var (
	blockchainInstance *Blockchain
	once               sync.Once
)

// GetBlockchain returns a singleton blockchain instance
func GetBlockchain() *Blockchain {
	once.Do(func() {
		blockchainInstance = CreateBlockchain()
	})
	return blockchainInstance
}

// GetAllBlocks returns all blocks from latest to genesis
func (bc *Blockchain) GetAllBlocks() []*Block {
	var blocks []*Block
	it := bc.Iterator()

	for {
		b := it.Next()
		blocks = append(blocks, b)
		if len(b.PrevBlockHash) == 0 {
			break
		}
	}
	return blocks
}

// GetBlocks is an alias for GetAllBlocks to keep server API clean
func (bc *Blockchain) GetBlocks() []*Block {
	return bc.GetAllBlocks()
}
