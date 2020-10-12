package assist

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
)

const randomSymbols = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const randomNumbers = "01234567890123456789012345678901234567890123456789012345678912"
const (
	symbolsIdxBits = 6                     // 6 bits to represent a letter index
	symbolsIdxMask = 1<<symbolsIdxBits - 1 // All 1-bits, as many as letterIdxBits
	symbolsIdxMax  = 63 / symbolsIdxBits   // # of letter indices fitting in 63 bits
)

var (
	randomSource  = rand.NewSource(time.Now().UnixNano())
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func RandomString(length int) string {
	return randomString(length, randomSymbols)
}

func RandomNumberString(length int) string {
	return randomString(length, randomNumbers)
}

func randomString(length int, source string) string {
	sb := strings.Builder{}
	sb.Grow(length)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := length-1, randomSource.Int63(), symbolsIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = randomSource.Int63(), symbolsIdxMax
		}
		if idx := int(cache & symbolsIdxMask); idx < len(source) {
			sb.WriteByte(source[idx])
			i--
		}
		cache >>= symbolsIdxBits
		remain--
	}
	return sb.String()
}

