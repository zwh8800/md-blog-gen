package util

import (
	"crypto/md5"
	"encoding/hex"
)

func MD5(src string) string {
	sumData := md5.Sum([]byte(src))
	return hex.EncodeToString(sumData[:])
}
