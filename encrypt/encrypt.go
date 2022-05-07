package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

type Encrypt struct {
	cb cipher.Block
	iv []byte
}

func New(key string, iv []byte) (*Encrypt, error) {
	cb, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	e := Encrypt{
		cb: cb,
		iv: iv,
	}
	return &e, nil
}

func (e *Encrypt) Encrypt(v string) string {
	padded := e.padding([]byte(v))
	encrypted := make([]byte, len(padded))
	enc := cipher.NewCBCEncrypter(e.cb, e.iv)
	enc.CryptBlocks(encrypted, padded)
	return base64.StdEncoding.EncodeToString(encrypted)
}

func (e *Encrypt) Decrypt(v string) (*string, error) {
	data, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		return nil, err
	}
	decrypted := make([]byte, len(data))
	dec := cipher.NewCBCDecrypter(e.cb, e.iv)
	dec.CryptBlocks(decrypted, data)
	d := string(e.unPadding(decrypted))
	return &d, nil
}

func (e *Encrypt) padding(data []byte) []byte {
	length := aes.BlockSize - (len(data) % aes.BlockSize)
	trailing := bytes.Repeat([]byte{byte(length)}, length)
	return append(data, trailing...)
}

func (e *Encrypt) unPadding(data []byte) []byte {
	return data[:len(data)-int(data[len(data)-1])]
}
