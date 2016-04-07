package util

import "testing"

func TestPngGen(t *testing.T) {
	const content = "https://lengzzz.com/note/4"

	_png(content, t)
	_jpeg(content, t)
}

func _png(content string, t *testing.T) {
	if err := GenerateQrcodePngFile(content, "/tmp/qr.png"); err != nil {
		t.Fatal(err)
	}
	base64String, err := GenerateQrcodePngBase64(content)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("base64:", base64String)

	dataUrl, err := GenerateQrcodePngDataUrl(content)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("dataUrl:", dataUrl)
}

func _jpeg(content string, t *testing.T) {
	if err := GenerateQrcodeJpegFile(content, "/tmp/qr.jpeg"); err != nil {
		t.Fatal(err)
	}
	base64String, err := GenerateQrcodeJpegBase64(content)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("base64:", base64String)

	dataUrl, err := GenerateQrcodeJpegDataUrl(content)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("dataUrl:", dataUrl)
}
