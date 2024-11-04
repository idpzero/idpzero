package idp

import (
	"fmt"

	"github.com/idpzero/idpzero/pkg/configuration"
	"github.com/idpzero/idpzero/pkg/validation"
)

func PrintValidation(config *configuration.ServerConfig) bool {

	var valid = true

	val := validation.NewValidationSet()
	for _, client := range config.Clients {
		set := validation.NewChecklist(fmt.Sprintf("Client '%s'", client.ID))
		_, validationErrors := NewClient(client)
		set.AddMany(validationErrors)

		if len(validationErrors) > 0 {
			valid = false
			val.AddChild(set)
		}
	}

	val.Render()
	return valid
}
