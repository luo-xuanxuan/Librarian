package utils

func Difference[T comparable](a, b []T) []T {
	var diff []T
	for _, elem := range a {
		if !Contains(b, elem) {
			diff = append(diff, elem)
		}
	}
	return diff
}
