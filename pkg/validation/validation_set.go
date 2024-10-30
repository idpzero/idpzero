package validation

var (
	_ Validation = &ValidationSet{}
)

func NewValidationSet() *ValidationSet {
	return &ValidationSet{
		validationWithChildren: validationWithChildren{
			children: []Validation{},
		},
	}
}

type ValidationSet struct {
	validationWithChildren
}

func (v *ValidationSet) Render() {
	v.render("")
}

// render implements Validation.
func (v *ValidationSet) render(indent string) {
	renderWithIndent(indent, v.validationWithChildren.children)
}
