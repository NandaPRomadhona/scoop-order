package src

// FindCountry checks if a string is present in a slice
func FindCountry(elements []string, element string) bool {
	for _, v := range elements {
		if string(v) == element {
			return true
		}
	}
	return false
}
