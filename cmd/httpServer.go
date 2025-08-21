package cmd

import (
	"github.com/Shubham0699/go-mini-blockchain/block"
	"github.com/Shubham0699/go-mini-blockchain/server"
	"github.com/spf13/cobra"
)

var port string

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Start HTTP server for blockchain",
	Run: func(cmd *cobra.Command, args []string) {
		bc := block.GetBlockchain()
		s := server.NewServer(bc)
		s.Start(port)
	},
}

func init() {
	httpCmd.Flags().StringVarP(&port, "port", "p", "8080", "Port to run HTTP server")
	rootCmd.AddCommand(httpCmd)
}
