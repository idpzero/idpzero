package console

import (
	"fmt"

	"github.com/fatih/color"
)

type Icon int

const (
	IconCheck Icon = iota
	IconCross
	IconDash
	IconQuestion
)

func resolveIcon(icon Icon) (string, func(fmt string, args ...interface{}) string) {
	switch icon {
	case IconCheck:
		return "✔", color.GreenString
	case IconCross:
		return "✖", color.RedString
	case IconDash:
		return "-", color.BlueString
	case IconQuestion:
		return "?", color.YellowString
	default:
		return "", color.WhiteString
	}
}

func PrintCheck(icon Icon, format string, args ...interface{}) {
	i, render := resolveIcon(icon)

	fmt.Printf("%s %s\n", render(i), color.WhiteString(format, args...))
}
