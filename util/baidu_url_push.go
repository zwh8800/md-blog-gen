package util

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/zwh8800/md-blog-gen/conf"
)

func PushUrlToBaidu(urls []string) ([]byte, error) {
	if !conf.Conf.Env.Prod {
		return nil, nil
	}
	buffer := &bytes.Buffer{}
	for _, url := range urls {
		buffer.WriteString(url + "\r\n")
	}

	resp, err := http.Post(conf.Conf.UrlPush.Baidu, "text/plain", buffer)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
