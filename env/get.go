package env

import (
	"os"
	"strconv"
	"strings"
	"time"
)

func GetStr(key string, def string) string {
	if v := getEnv(key); len(v) > 0 {
		return v
	}
	return def
}

func GetInt(key string, def int) int {
	i, err := strconv.Atoi(getEnv(key))
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

func GetDuration(key string, def time.Duration) time.Duration {
	if d, _ := time.ParseDuration(getEnv(key)); d > 0 {
		return d
	}
	return def
}

func getEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key));
}
