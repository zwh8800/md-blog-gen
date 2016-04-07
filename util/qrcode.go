package util

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func GenerateQrcodeImage(content string) (image.Image, error) {
	qrcode, err := qr.Encode(content, qr.M, qr.Auto)
	if err != nil {
		return nil, err
	}
	return barcode.Scale(qrcode, 100, 100)
}

func GenerateQrcodePngFile(content string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	img, err := GenerateQrcodeImage(content)
	if err != nil {
		return err
	}
	return png.Encode(file, img)
}

func GenerateQrcodePngBase64(content string) (string, error) {
	img, err := GenerateQrcodeImage(content)
	if err != nil {
		return "", err
	}
	output := &bytes.Buffer{}
	base64Encoder := base64.NewEncoder(base64.StdEncoding, output)

	if err := png.Encode(base64Encoder, img); err != nil {
		return "", err
	}

	return output.String(), nil
}

func GenerateQrcodePngDataUrl(content string) (string, error) {
	base64String, err := GenerateQrcodePngBase64(content)
	if err != nil {
		return "", err
	}
	return "data:" + "image/png;base64," + base64String, nil
}

func GenerateQrcodeJpegFile(content string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	img, err := GenerateQrcodeImage(content)
	if err != nil {
		return err
	}
	return jpeg.Encode(file, img, &jpeg.Options{jpeg.DefaultQuality})
}

func GenerateQrcodeJpegBase64(content string) (string, error) {
	img, err := GenerateQrcodeImage(content)
	if err != nil {
		return "", err
	}
	output := &bytes.Buffer{}
	base64Encoder := base64.NewEncoder(base64.StdEncoding, output)

	if err := jpeg.Encode(base64Encoder, img, &jpeg.Options{jpeg.DefaultQuality}); err != nil {
		return "", err
	}

	return output.String(), nil
}

func GenerateQrcodeJpegDataUrl(content string) (string, error) {
	base64String, err := GenerateQrcodeJpegBase64(content)
	if err != nil {
		return "", err
	}
	return "data:" + "image/jpeg;base64," + base64String, nil
}
