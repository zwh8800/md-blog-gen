package util

import (
	"bytes"
	"net/http"

	"github.com/zwh8800/md-blog-gen/conf"
)

func PushUrlToBaidu(urls []string) error {
	buffer := &bytes.Buffer{}
	for _, url := range urls {
		buffer.WriteString(url)
	}

	if _, err := http.Post(conf.Conf.UrlPush.Baidu, "text/plain", buffer); err != nil {
		return err
	}

	return nil
}
