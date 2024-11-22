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
		return "(✔)", color.GreenString
	case IconCross:
		return "(✖)", color.RedString
	case IconDash:
		return "(-)", color.WhiteString
	case IconQuestion:
		return "(?)", color.BlueString
	default:
		return "", color.WhiteString
	}
}

type Check struct {
	format string
	args   []interface{}
}

func NewCheck(format string, args ...interface{}) Check {
	return Check{
		format: format,
		args:   args,
	}
}

func (c Check) Print(icon Icon, format string, args ...interface{}) {
	i, render := resolveIcon(icon)

	fmt.Printf("%s %s %s\n", render(i), color.WhiteString(c.format, c.args...), render("["+format+"]", args...))
}

func PrintCheck(icon Icon, format string, args ...interface{}) {
	i, render := resolveIcon(icon)

	fmt.Printf("%s %s\n", render(i), color.WhiteString(format, args...))
}
