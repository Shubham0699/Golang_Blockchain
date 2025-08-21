package cmd

import (
	"fmt"
	"github.com/Shubham0699/go-mini-blockchain/block"
	"github.com/spf13/cobra"
)

var printChainCmd = &cobra.Command{
	Use:   "printchain",
	Short: "Print all blocks in the blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		bc := block.GetBlockchain()
		defer bc.Close() // important to close DB after use

		it := bc.Iterator()

		for {
			blk := it.Next()

			fmt.Printf("\nHash: %x\nPrevHash: %x\nData: %s\n\n",
				blk.Hash, blk.PrevBlockHash, blk.Data)

			// stop when we reach the genesis block
			if len(blk.PrevBlockHash) == 0 {
				break
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(printChainCmd)
}
