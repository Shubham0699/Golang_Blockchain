package cmd

import (
	"fmt"
	"github.com/Shubham0699/go-mini-blockchain/block"

	"github.com/spf13/cobra"
)

var data string

var addBlockCmd = &cobra.Command{
	Use:   "addblock",
	Short: "Add a block to the blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		bc := block.GetBlockchain()
		bc.AddBlock(data)
		fmt.Println("✅ Block added with data:", data)
	},
}

func init() {
	addBlockCmd.Flags().StringVarP(&data, "data", "d", "", "Block data")
	addBlockCmd.MarkFlagRequired("data")
	rootCmd.AddCommand(addBlockCmd)
}
