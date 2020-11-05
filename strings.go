package assist

import (
	"net/url"
)

func IsStringUrl(v string) bool {
	url, err := url.ParseRequestURI(v)
	if err != nil {
		return false
	}
	return url != nil && url.IsAbs()
}

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
