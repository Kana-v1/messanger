package collections

func Contains(value interface{}, arr ...interface{}) interface{} {
	for _, el := range arr {
		if value == el {
			return value
		}
	}

	return nil
}

func Remove(value interface{}, arr ...interface{}) bool {
	for i, el := range arr {
		if value == el {
			newArr := arr[0:i]
			newArr = append(newArr, arr[i+1:])
			arr = newArr
			return true
		}
	}
	return false
}

func RemoveIndex(s []interface{}, index int) []interface{} {
	return append(s[:index], s[index+1:]...)
}
