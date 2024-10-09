package util

import (
	crand "crypto/rand"
	"math/big"
	"strconv"
)

func ParseInt(v string) int {
	x, e := strconv.ParseInt(v, 10, 32)
	if e != nil {
		return 0
	}
	return int(x)
}

const charsets = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func BuildSN(length int) string {
	result := make([]byte, length)
	fixedSize := int64(len(charsets))
	for i := range result {
		num, err := crand.Int(crand.Reader, big.NewInt(fixedSize))
		if err != nil {
			num = big.NewInt(0)
		}
		result[i] = charsets[num.Int64()]
	}
	return string(result)
}
