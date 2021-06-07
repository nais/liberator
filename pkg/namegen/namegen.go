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
// Can be used to avoid collisions in the Kubernetes namespace.
//
// e.g. ShortName("foobarbaz", 16) --> "foobarb-12345678"
func ShortName(basename string, maxlen int) (string, error) {
	maxlen -= SuffixLength
	hasher := crc32.NewIEEE()
	_, err := hasher.Write([]byte(basename))
	if err != nil {
		return "", err
	}
	hashStr := fmt.Sprintf("%x", hasher.Sum32())

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

func formatName(basename, suffix string, maxlen int) string {
	if len(basename) > maxlen {
		basename = basename[:maxlen]
	}
	return fmt.Sprintf("%s-%s", basename, suffix)
}
