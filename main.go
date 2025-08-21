package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Shubham0699/go-mini-blockchain/block"
	"github.com/Shubham0699/go-mini-blockchain/p2p"
	"github.com/Shubham0699/go-mini-blockchain/tx"
	"github.com/Shubham0699/go-mini-blockchain/wallet"
)

func main() {
	// Load blockchain
	bc := block.GetBlockchain()
	defer bc.Close()

	// Start P2P node
	node := p2p.NewNode("localhost:3000", bc)
	go node.StartServer()

	// Automatic mining every 10 seconds
	go func() {
		for {
			w, _ := wallet.NewWallet()
			cbTx := tx.NewCoinbaseTX(w.Address(), 50)
			node.Blockchain.MineBlock([]*tx.Transaction{cbTx})

			blocks := node.Blockchain.GetBlocks()
			node.BroadcastBlock(blocks[0])

			log.Println("✅ Auto-mined block to:", w.Address())
			time.Sleep(10 * time.Second)
		}
	}()

	// CLI for manual commands
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nCommands: mine <address> | print | exit")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		args := strings.Split(input, " ")

		switch args[0] {
		case "mine":
			if len(args) < 2 {
				fmt.Println("Usage: mine <address>")
				continue
			}
			address := args[1]
			cbTx := tx.NewCoinbaseTX(address, 50)
			node.Blockchain.MineBlock([]*tx.Transaction{cbTx})

			blocks := node.Blockchain.GetBlocks()
			node.BroadcastBlock(blocks[0])

			fmt.Println("✅ Mined block to:", address)

		case "print":
			printBlockchain(node.Blockchain)

		case "exit":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Unknown command")
		}
	}
}

func printBlockchain(bc *block.Blockchain) {
	fmt.Println("\n📜 Blockchain History:")
	blocks := bc.GetBlocks()
	for _, b := range blocks {
		fmt.Printf("---------------------------\n")
		fmt.Printf("Hash: %x\n", b.Hash)
		fmt.Printf("PrevHash: %x\n", b.PrevBlockHash)
		fmt.Printf("Nonce: %d\n", b.Nonce)
		fmt.Printf("Timestamp: %d\n", b.Timestamp)
		if len(b.Transactions) > 0 {
			fmt.Println("Transactions:")
			for _, t := range b.Transactions {
				fmt.Printf(" - TXID: %x\n", t.ID)
				fmt.Printf("   Inputs: %+v\n", t.Vin)
				fmt.Printf("   Outputs: %+v\n", t.Vout)
			}
		} else {
			fmt.Printf("Data: %s\n", string(b.Data))
		}
	}
	fmt.Println("---------------------------")
}
