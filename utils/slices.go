package utils

func Contains[T comparable](slice []T, elem T) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}
	return false
}

func Difference[T comparable](a, b []T) []T {
	var diff []T
	for _, elem := range a {
		if !Contains(b, elem) {
			diff = append(diff, elem)
		}
	}
	return diff
}

func Intersection[T comparable](a, b []T) []T {
	var inter []T
	// Use a map to track seen elements for improved efficiency
	seen := make(map[T]bool)
	for _, item := range b {
		seen[item] = true
	}

	for _, item := range a {
		if _, ok := seen[item]; ok && Contains(b, item) {
			inter = append(inter, item)
		}
	}
	return inter
}
