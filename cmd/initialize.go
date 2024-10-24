package cmd

import (
	"fmt"

	"math/rand"

	"github.com/spf13/cobra"
)

var initializeCmd = &cobra.Command{
	Use:   "init",
	Short: "Initalize configuration for idpzero",
	Long:  `Setup the configuration and data directory for idpzero`,
	RunE: func(cmd *cobra.Command, args []string) error {

		fmt.Println(*location)

		fmt.Println(randSeq(64))
		iss := fmt.Sprintf("http://localhost:%d", 4379)
		fmt.Println(iss)
		return nil
	},
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
