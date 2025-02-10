package stringutil

import (
	"fmt"
	"hash/crc32"
	"math/rand"
	"time"
)

const strTrimMiddleTruncate = "---[truncated]---"
const strTrimRightTruncate = "..."

func StrTrimRight(s string, maxlen int) string {
	l := len(strTrimRightTruncate)
	if maxlen > l {
		return s[:maxlen-l] + strTrimRightTruncate
	}
	return s[:maxlen]
}

func StrTrimMiddle(s string, maxlen int) string {
	if len(s) <= maxlen {
		return s
	}
	newlen := maxlen - len(strTrimMiddleTruncate)
	if newlen < len(strTrimMiddleTruncate) {
		return StrTrimRight(s, maxlen)
	}
	partlen := int(float64(newlen) / 2)
	return s[:partlen] + strTrimMiddleTruncate + s[len(s)-partlen:]
}

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func CRC32(input string) string {
	hasher := crc32.NewIEEE()
	// crc32 hasher always returns nil error
	_, _ = hasher.Write([]byte(input))
	return fmt.Sprintf("%x", hasher.Sum32())
}

// Procedurally generate a short string with hash that can be calculated using the base name
func UniqueWithHash(basename string, maxlen int) string {
	maxlen -= 9 // 8 bytes of hexadecimal hash and 1 byte of separator
	hashStr := CRC32(basename)
	if len(basename) > maxlen {
		basename = basename[:maxlen]
	}
	return basename + "-" + hashStr
}
