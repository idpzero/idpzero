package pages

import "strconv"

// formatInt renders an integer as a string for use in templates.
func formatInt(i int) string {
	return strconv.Itoa(i)
}
