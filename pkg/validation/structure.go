package validation

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	_ Validation = &checkedItem{}
	_ Validation = &validationGroup{}
)

type Validation interface {
	Render()
	AddChild(item Validation)
}

type checkedItem struct {
	passed   bool
	children []Validation
	title    string
	err      string
}

func (v *checkedItem) Render() {
	var mark = color.RedString(" x ")
	if v.passed {
		mark = color.GreenString(" ✓ ")
	}

	fmt.Println(mark, v.title)
	if v.err != "" {
		color.Red("   ", v.err)
	}

	for _, child := range v.children {
		child.Render()
	}

}

func (v *checkedItem) AddChild(item Validation) {
	v.children = append(v.children, item)
}

type validationGroup struct {
	children []Validation
	title    string
}

func (v *validationGroup) Render() {
	fmt.Println(v.title)

	for _, child := range v.children {
		child.Render()
	}

	fmt.Println()
}

func (v *validationGroup) AddChild(item Validation) {
	v.children = append(v.children, item)
}

func NewValidation(title string) Validation {
	return &validationGroup{
		title: title,
	}
}

func NewCheckedValidation(passed bool, title string) Validation {
	return &checkedItem{
		title:  title,
		passed: passed,
	}
}
