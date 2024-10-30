package validation

type Validation interface {
	render(indent string)
}

type ValidationWithChildren interface {
	AddChild(item Validation)
}

type validationWithChildren struct {
	children []Validation
}

// AddChild adds a child validation to the set
func (vs *validationWithChildren) AddChild(child Validation) {
	vs.children = append(vs.children, child)
}

func renderWithIndent(indent string, v []Validation) {
	for _, child := range v {
		child.render(indent)
	}
}
