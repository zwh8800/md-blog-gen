package sitemap

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestSitemap(t *testing.T) {
	si := New()
	for i := 0; i < 10; i++ {
		si.Add(
			&Item{
				Link:    fmt.Sprintf("http://localhost/item/%d", i),
				Updated: time.Now()})
	}
	if err := Xml(si, os.Stdout); err != nil {
		t.Errorf("Error on write xml: %v", err)
	}
	os.Stdout.Write([]byte("\n"))
}
