package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type Encript struct {
	cb cipher.Block
	iv []byte
}

func New(key string, iv []byte) (*Encript, error) {
	cb, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	e := Encript{
		cb: cb,
		iv: iv,
	}
	return &e, nil
}

func (e *Encript) Encrypt(v string) string {
	padded := e.padding([]byte(v))
	encrypted := make([]byte, len(padded))
	enc := cipher.NewCBCEncrypter(e.cb, e.iv)
	enc.CryptBlocks(encrypted, padded)
	return base64.StdEncoding.EncodeToString(encrypted)
}

func (e *Encript) Decrypt(v string) string {
	data, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		panic(err)
	}
	decrypted := make([]byte, len(data))
	dec := cipher.NewCBCDecrypter(e.cb, e.iv)
	dec.CryptBlocks(decrypted, data)
	return string(e.unpadding(decrypted))
}

func (e *Encript) padding(data []byte) []byte {
	length := aes.BlockSize - (len(data) % aes.BlockSize)
	trailing := bytes.Repeat([]byte{byte(length)}, length)
	return append(data, trailing...)
}

func (e *Encript) unpadding(data []byte) []byte {
	return data[:len(data)-int(data[len(data)-1])]
}
