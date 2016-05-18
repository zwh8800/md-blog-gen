package util

import (
	"bytes"

	"github.com/mozillazg/go-pinyin"
)

func Pinyin(title string) string {
	py := pinyin.Pinyin(title, pinyin.NewArgs())
	sb := &bytes.Buffer{}
	for _, p := range py {
		sb.WriteString(p[0])
		sb.WriteRune(' ')
	}
	return sb.String()
}
