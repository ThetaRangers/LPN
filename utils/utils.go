package utils

func Contains(slice []string, string string) bool {
	for _, x := range slice {
		if x == string {
			return true
		}
	}

	return false
}

func RemoveFromList(slice []string, target string) []string {
	var list []string

	for i, x := range slice {
		if x == target {
			list = append(slice[:i], slice[i+1:]...)
			return list
		}
	}

	return slice
}