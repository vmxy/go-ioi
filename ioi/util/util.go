package util

import "strconv"

func ParseInt(v string) int {
	x, e := strconv.ParseInt(v, 10, 32)
	if e != nil {
		return 0
	}
	return int(x)
}
