package initialize

import (
	"fmt"
	"os"
	"path"

	"math/rand"

	"github.com/spf13/cobra"
)

func Register(parent *cobra.Command) {

	var location *string
	var port *int

	var cmd = &cobra.Command{
		Use:   "init",
		Short: "Initalize configuration for idpzero",
		Long:  `Setup the configuration and data directory for idpzero`,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			location = cmd.Flags().StringP("location", "l", path.Join(cwd, ".idpzero"), "Location to store configuration")
			port = cmd.Flags().IntP("port", "p", 4379, "Port to serve the IDP server on")

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Println(*location)
			fmt.Println(*port)

			fmt.Println(randSeq(64))
			iss := fmt.Sprintf("http://localhost:%d", *port)
			fmt.Println(iss)
			return nil
		},
	}

	parent.AddCommand(cmd)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
