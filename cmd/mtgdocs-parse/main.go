package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Manbeardo/mtgdocs.info/parse"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "mtgdocs-parse",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "cr [filename]",
		Short: "parses the comprehensive rules",
		RunE:  ParseCR,
	})

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
}

func ParseCR(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid number of args")
	}
	file, err := os.Open(args[0])
	if err != nil {
		return err
	}
	cr := parse.ParseCR(file)
	json.NewEncoder(os.Stdout).Encode(cr["HEAD"].Children)
	return nil
}
