package namegen

import (
	"fmt"
	"hash/crc32"
)

// Append the string's CRC32 hash to the string and truncate it to a maximum length.
// Can be used to avoid collisions in the Kubernetes namespace.
//
// e.g. ShortName("foobarbaz", 16) --> "foobarb-12345678"
func ShortName(basename string, maxlen int) (string, error) {
	maxlen -= 9 // 8 bytes of hexadecimal hash and 1 byte of separator
	hasher := crc32.NewIEEE()
	_, err := hasher.Write([]byte(basename))
	if err != nil {
		return "", err
	}
	hashStr := fmt.Sprintf("%x", hasher.Sum32())
	if len(basename) > maxlen {
		basename = basename[:maxlen]
	}
	return basename + "-" + hashStr, nil
}
