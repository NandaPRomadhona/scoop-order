package src

// FindString checks if a string is present in a slice
func FindString(elements []string, element string) bool {
	for _, v := range elements {
		if string(v) == element {
			return true
		}
	}
	return false
}

func FindInt(what int32, where []int32) (idx bool) {
	for _, v := range where {
		if v == what {
			return true
		}
	}
	return false
}