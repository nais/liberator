package strings

import "slices"

// Helper functions to check and remove string from a slice of strings.
func ContainsString(slice []string, s string) bool {
	return slices.Contains(slice, s)
}

func RemoveString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
