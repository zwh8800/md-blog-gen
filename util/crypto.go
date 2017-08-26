package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io"
)

func RsaEncrypt(src []byte, key []byte) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("public key error")
	}
	itf, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := itf.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, src)
}

func Base64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func DecryptAesCBC(msg string, key []byte, iv []byte) ([]byte, error) {
	cipherText, _ := hex.DecodeString(msg)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(cipherText) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	decrypted := make([]byte, len(cipherText))
	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(decrypted, cipherText)
	decrypted = PKCS5UnPadding(decrypted)
	return decrypted, nil
}

// PKCS5UnPadding unpads the given data and returns the original data without
// PKCS5 padding
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// remove the last byte
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

func AesEncrypt(msg string, key []byte) (string, error) {
	plaintext := []byte(msg)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the cipherText.
	cipherText := make([]byte, aes.BlockSize+len(plaintext))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", errors.New("产生随机iv出错")
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], plaintext)
	// convert to base64
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func AesDecrypt(msg string, key []byte) (string, error) {
	cipherText, _ := base64.URLEncoding.DecodeString(msg)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.New("密钥不符合要求")
	}
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the cipherText.
	if len(cipherText) < aes.BlockSize {
		return "", errors.New("加密数据太短")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	// XORKeyStream can work in-place if the two arguments are the same.
	stream.XORKeyStream(cipherText, cipherText)
	return string(cipherText), nil
}
