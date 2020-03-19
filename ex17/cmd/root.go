package cmd

import (
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var encodingkey string

// RootCmd is the root commnad.
var RootCmd = &cobra.Command{
	Use:   "secret",
	Short: "Secret is an API key and other secrets manager",
}

func secretsPath() string {
	home, _ := homedir.Dir()
	return filepath.Join(home, ".secretes")
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&encodingkey, "key", "k", "", "The key to use when encoding and decoding secrets.")
}
