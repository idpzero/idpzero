package idp

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/idpzero/idpzero/pkg/configuration"
)

func PrintValidation(config *configuration.IDPConfiguration) bool {

	var valid = true

	fmt.Printf("Validaing %d clients:\n", len(config.Clients))

	for _, client := range config.Clients {
		c, validationErrors := NewClient(client)

		configuration.PrintCheck(c.IsValid(), fmt.Sprintf("Client '%s'", client.ID))

		if !c.IsValid() {
			valid = false // make sure we wont apply the config
			for _, err := range validationErrors {
				color.Red("     " + err.Error())
			}
		}
	}
	fmt.Println()
	return valid
}
