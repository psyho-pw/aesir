package utils

func MapValues[T interface{}](obj map[string]T) []T {
	arr := make([]T, 0)

	for _, value := range obj {
		arr = append(arr, value)
	}

	return arr
}

func MapKeys(obj map[string]interface{}) []string {
	arr := make([]string, 0)

	for key, _ := range obj {
		arr = append(arr, key)
	}

	return arr
}
