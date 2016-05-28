package util

import "testing"

func TestAutoClearKVCache(t *testing.T) {
	cache := NewAutoClearKVCache(10)
	cache.Put("haha", "+1s")

	for i := 0; i < 10; i++ {
		if value, ok := cache.Get("haha"); ok {
			if value != "+1s" {
				t.Errorf("get value not ok")
				t.Fail()
			}
		} else {
			t.Errorf("get not ok")
			t.Fail()
		}
	}
	if _, ok := cache.Get("haha"); ok {
		t.Errorf("cache not clear")
		t.Fail()
	}
}
