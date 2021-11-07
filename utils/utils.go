package utils

func Contains(slice []string, string string) bool {
	for _, x := range slice {
		if x == string {
			return true
		}
	}

	return false
}