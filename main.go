package main

import (
    "fmt"
    "github.com/Shubham0699/go-mini-blockchain/block" // ✅ make sure import path matches
)

func main() {
    bc := block.NewBlockchain() // ✅ Now uses correct package prefix

    // Add some blocks
    bc.AddBlock("Alice sent 5 coins to Bob")
    bc.AddBlock("Bob sent 2 coins to Charlie")
    bc.AddBlock("Charlie sent 1 coin to Dave")

    // Print the blockchain
    for i, blk := range bc.Blocks {
        fmt.Println("====================================")
        fmt.Printf("🧱 Block #%d\n", i)
        fmt.Printf("🕒 Timestamp: %d\n", blk.Timestamp)
        fmt.Printf("📜 Data: %s\n", blk.Data)
        fmt.Printf("🔗 Prev. Hash: %x\n", blk.PrevBlockHash)
        fmt.Printf("🧬 Hash: %x\n", blk.Hash)
    }

    // ✅ Validate the entire blockchain at the end
    fmt.Println("====================================")
    bc.IsValid()
}
