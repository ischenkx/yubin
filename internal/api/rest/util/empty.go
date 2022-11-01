package util

func ValidateEmptySlice[T any](s []T) []T {
	if s == nil {
		return []T{}
	}
	return s
}

func ValidateEmptyMap[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return map[K]V{}
	}
	return m
}
