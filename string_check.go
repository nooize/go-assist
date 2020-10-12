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
