package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var helloCmd = &cobra.Command{
	Use:   "hello",
	Short: "Just a hello world command",

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello world!")
	},
}
