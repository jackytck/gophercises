package cmd

import (
	"fmt"

	"github.com/jackytck/gophercises/ex17/vault"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.File(encodingkey, secretsPath())
		key := args[0]
		val, err := v.Get(key)
		if err != nil {
			fmt.Println("No value set!")
			return
		}
		fmt.Printf("%s = %s\n", key, val)
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
