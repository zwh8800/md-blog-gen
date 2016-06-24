package util

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
)

// 62个字符, 需要6bit做索引(2 ^ 6 = 64)
var charTable = [...]rune{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k',
	'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v',
	'w', 'x', 'y', 'z', '0', '1', '2', '3', '4', '5', '6',
	'7', '8', '9', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
	'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S',
	'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
}

func ShortenUrl(url string) []string {
	shortUrlList := make([]string, 0, 4)

	sumData := md5.Sum([]byte(url))
	for i := 0; i < 4; i++ {
		part := sumData[i*4 : i*4+4]
		partUint := binary.BigEndian.Uint32(part)
		partUint &= 0x3fffffff

		shortUrlBuffer := &bytes.Buffer{}
		// 将30bit分成6份, 每份5bit
		for j := 0; j < 6; j++ {
			index := partUint % 62

			shortUrlBuffer.WriteRune(charTable[index])
			partUint = partUint >> 5
		}
		shortUrlList = append(shortUrlList, shortUrlBuffer.String())
	}
	return shortUrlList
}
