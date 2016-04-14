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
	// 把md5sum分成4份, 每份4个字节
	for i := 0; i < 4; i++ {
		part := sumData[i*4 : i*4+4]
		// 将4字节当作一个整数
		partUint := binary.BigEndian.Uint32(part)
		// 只取整数的后30个bit, 用&屏蔽掉高位
		partUint &= 0x3fffffff

		shortUrlBuffer := &bytes.Buffer{}
		// 将30bit分成6份, 每份5bit
		for j := 0; j < 6; j++ {
			// 0x3d = 0b00111101, 多取一位, 扔自己一位
			// 巧妙之处, 解决了5bit不够索引62个字符的尴尬
			index := partUint & 0x3d

			shortUrlBuffer.WriteRune(charTable[index])
			partUint = partUint >> 5
		}
		shortUrlList = append(shortUrlList, shortUrlBuffer.String())
	}
	return shortUrlList
}
