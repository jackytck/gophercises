package cmd

import (
	"fmt"

	"github.com/jackytck/gophercises/ex17/vault"
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Sets a secret in your secret storage",
	Run: func(cmd *cobra.Command, args []string) {
		v := vault.File(encodingkey, secretsPath())
		key, val := args[0], args[1]
		err := v.Set(key, val)
		if err != nil {
			panic(err)
		}
		fmt.Println("Value set successfully!")
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
