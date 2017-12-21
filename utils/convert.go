package utils

import (
	"strconv"
)

func Int642String(val int64) string {
	return strconv.FormatInt(val, 10)
}

func String2Int64(val string) (int64, error) {
	return strconv.ParseInt(val, 10, 64)
}
