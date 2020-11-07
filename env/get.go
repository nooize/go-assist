package env

import (
	"os"
	"strconv"
	"strings"
)

func GetStr(key string, def string) string {
	if v := strings.TrimSpace(os.Getenv(key)); len(v) > 0 {
		return v
	}
	return def
}

func GetInt(key string, def int) int {
	i, err := strconv.Atoi(strings.TrimSpace(os.Getenv(key)))
	if err != nil {
		return def
	}
	return i
}

func GetPositiveInt(key string, def int) int {
	i :=  GetInt(key, def)
	if i < 0 {
		return def
	}
	return i
}
