package utils

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
