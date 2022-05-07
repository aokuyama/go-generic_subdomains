package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
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
	if len(decrypted)%dec.BlockSize() != 0 {
		return nil, errors.New("crypto/cipher: input not full blocks")
	}
	if len(decrypted) < len(data) {
		return nil, errors.New("crypto/cipher: output smaller than input")
	}
	dec.CryptBlocks(decrypted, data)
	d := string(e.unPadding(decrypted))
	return &d, nil
}

func (e *Encrypt) DecryptSafe(v string) string {
	s, err := e.Decrypt(v)
	if err != nil {
		return v
	}
	return *s
}

func (e *Encrypt) padding(data []byte) []byte {
	length := aes.BlockSize - (len(data) % aes.BlockSize)
	trailing := bytes.Repeat([]byte{byte(length)}, length)
	return append(data, trailing...)
}

func (e *Encrypt) unPadding(data []byte) []byte {
	return data[:len(data)-int(data[len(data)-1])]
}
