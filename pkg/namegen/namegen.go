package namegen

import (
	"fmt"
	"hash/crc32"

	"github.com/nais/liberator/pkg/keygen"
)

const (
	SuffixLength = 9 // 8 bytes of hexadecimal hash/random characters and 1 byte of separator
)

// Append the string's CRC32 hash to the string and truncate it to a maximum length.
// Can be used to avoid collisions in the Kubernetes namespace. CRC is deterministic.
//
// e.g. ShortName("foobarbaz", 16) --> "foobarb-12345678", which is longer than the original name
func ShortName(basename string, maxlen int) (string, error) {
	maxlen -= SuffixLength
	hasher := crc32.NewIEEE()
	_, err := hasher.Write([]byte(basename))
	if err != nil {
		return "", err
	}
	hashStr := fmt.Sprintf("%08x", hasher.Sum32())

	return formatName(basename, hashStr, maxlen), nil
}

// Suffix a string with random letters and numbers, and truncate it to a maximum length.
func RandShortName(basename string, maxlen int) string {
	maxlen -= SuffixLength
	suffix := keygen.RandStringBytes(SuffixLength - 1)

	return formatName(basename, suffix, maxlen)
}

func PrefixedRandShortName(prefix, basename string, maxlen int) string {
	return RandShortName(fmt.Sprintf("%s-%s", prefix, basename), maxlen)
}

// SuffixedShortName generates a ShortName and appends a given suffix to this, e.g. an incremental counter.
// The resulting string is truncated to the given maximum length.
func SuffixedShortName(basename, suffix string, maxlen int) (string, error) {
	maxlen -= len(suffix) + 1 // length of suffix + 1 byte of separator
	shortName, err := ShortName(basename, maxlen)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", shortName, suffix), nil
}

func formatName(basename, suffix string, maxlen int) string {
	if len(basename) > maxlen {
		basename = basename[:maxlen]
	}
	return fmt.Sprintf("%s-%s", basename, suffix)
}
