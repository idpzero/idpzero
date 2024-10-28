package configuration

import (
	"fmt"

	"github.com/fatih/color"
)

func PrintCheck(passed bool, msg string) {
	var mark = color.RedString(" x ")
	if passed {
		mark = color.GreenString(" ✓ ")
	}

	fmt.Println(mark, msg)
}
