package assist

func IsArrayContainString(list *[]string, key string) bool {
	if list != nil {
		for _, v := range *list {
			if v == key {
				return true
			}
		}
	}
	return false
}
