package validation

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	_ Validation = &Check{}
	_ Validation = &ValidationGroup{}
)

type Validation interface {
	render(indent int)
	AddChild(item Validation)
}

type Check struct {
	passed   bool
	children []Validation
	title    string
	err      string
	options  []string
}

func (v *Check) render(indent int) {
	var mark = color.RedString(" x ")
	if v.passed {
		mark = color.GreenString(" ✓ ")
	}

	fmt.Println(mark, v.title)
	if v.err != "" {
		color.Red("   ", v.err)
	}

	for _, child := range v.children {
		child.render(indent + 1)
	}

}

func (v *Check) AddChild(item Validation) {
	v.children = append(v.children, item)
}

func (v *Check) WithOptions(items []string) {
	v.options = items
}

func NewCheck(passed bool, title string) *Check {
	return &Check{
		title:   title,
		passed:  passed,
		options: []string{},
	}
}

type ValidationGroup struct {
	children []Validation
	title    string
}

func (v *ValidationGroup) Render() {
	v.render(0)
}

func (v *ValidationGroup) render(indent int) {
	fmt.Println(v.title)

	for _, child := range v.children {
		child.render(indent + 1)
	}

	fmt.Println()
}

func (v *ValidationGroup) AddChild(item Validation) {
	v.children = append(v.children, item)
}

func NewValidation(title string) *ValidationGroup {
	return &ValidationGroup{
		title: title,
	}
}
