package util

import "net/http"

func WriteContentType(w http.ResponseWriter, value string) {
	header := w.Header()

	if val := header.Get("Content-Type"); len(val) == 0 {
		header.Set("Content-Type", value)
	}
}
