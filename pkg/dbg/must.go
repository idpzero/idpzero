package dbg

func MustOrFalse(val bool, err error) bool {
	if err != nil {
		return false
	}

	return val
}
