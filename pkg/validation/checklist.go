package validation

import (
	"fmt"

	"github.com/fatih/color"
)

var (
	_ Validation = &Checklist{}
)

type ChecklistItem struct {
	passed  bool
	title   string
	err     error
	options []string
	val     interface{}
}

func (v *ChecklistItem) render(indent string) {
	var mark = color.RedString(" x ")
	if v.passed {
		mark = color.GreenString(" ✓ ")
	}

	value := ""

	if v.val != nil {
		value = color.MagentaString(" (%v)", v.val)
	}

	fmt.Println(indent, mark, v.title, value)

	if v.err != nil {
		color.RedString(indent+"   ", v.err)
	}

	if len(v.options) > 0 {
		fmt.Println(indent, color.YellowString("    expected value(s):"))
		for _, option := range v.options {
			fmt.Println(indent, color.YellowString("     - %s", option))
		}
	}

}

func (v *ChecklistItem) WithError(err error) *ChecklistItem {
	v.err = err

	return v
}

func (v *ChecklistItem) WithOptions(items []string) *ChecklistItem {
	v.options = items

	return v
}

func (v *ChecklistItem) WithValue(val interface{}) *ChecklistItem {
	v.val = val

	return v
}

func NewChecklistItem(passed bool, title string) *ChecklistItem {
	return &ChecklistItem{
		title:   title,
		passed:  passed,
		options: []string{},
	}
}

type Checklist struct {
	title string
	items []*ChecklistItem
}

func (v *Checklist) Add(item *ChecklistItem) {
	v.items = append(v.items, item)
}

func (v *Checklist) AddMany(item []*ChecklistItem) {
	v.items = append(v.items, item...)
}

func (v *Checklist) render(indent string) {
	fmt.Println(v.title)

	for _, child := range v.items {
		child.render(indent)
	}

	fmt.Println()
}

func NewChecklist(title string) *Checklist {
	return &Checklist{
		title: title,
		items: []*ChecklistItem{},
	}
}
