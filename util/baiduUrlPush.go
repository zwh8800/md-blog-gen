package util

import (
	"bytes"
	"net/http"

	"github.com/zwh8800/md-blog-gen/conf"
)

func PushUrlToBaidu(urls []string) error {
	if !conf.Conf.Env.Prod {
		return nil
	}
	buffer := &bytes.Buffer{}
	for _, url := range urls {
		buffer.WriteString(url + "\r\n")
	}

	if _, err := http.Post(conf.Conf.UrlPush.Baidu, "text/plain", buffer); err != nil {
		return err
	}

	return nil
}
