package cache

import (
	"bytes"
	"io"
	"io/ioutil"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"github.com/zwh8800/md-blog-gen/service"
	"github.com/zwh8800/md-blog-gen/util"
	redis "gopkg.in/redis.v5"
)

type readableResponseWriter struct {
	gin.ResponseWriter
	buffer *bytes.Buffer
}

func newReadableResponseWriter(writer gin.ResponseWriter) (*readableResponseWriter, io.Reader) {
	buffer := new(bytes.Buffer)
	return &readableResponseWriter{
		ResponseWriter: writer,
		buffer:         buffer,
	}, buffer
}

func (w *readableResponseWriter) Write(data []byte) (int, error) {
	if n, err := w.ResponseWriter.Write(data); err != nil {
		return n, err
	}
	return w.buffer.Write(data)
}

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.EscapedPath()
		data, err := service.FindCache(path)
		if err != nil {
			if err == redis.Nil {
				newWriter, reader := newReadableResponseWriter(c.Writer)
				c.Writer = newWriter
				defer func() {
					if c.Writer.Status() == http.StatusOK {
						data, _ := ioutil.ReadAll(reader)
						service.AddCache(path, string(data))
					}
				}()
			} else {
				glog.Error(err)
			}
			c.Next()
		} else {
			util.WriteContentType(c.Writer, []string{"text/html; charset=utf-8"})
			c.Writer.Write([]byte(data))
			c.Abort()
		}
	}
}
