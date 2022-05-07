package encrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncrypt(t *testing.T) {
	iv := []byte{0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}
	e, _ := New("BuiyLkrD5bY5PKsVX5wAEw7PdAJGBzys", iv)
	str := "あいうえお"
	enc := e.Encrypt(str)
	assert.NotEqual(t, enc, str)
	assert.NotEqual(t, enc, e.Encrypt("あいうえお1"))
	assert.Equal(t, enc, e.Encrypt(str), "何度暗号化しても同じ結果")
	assert.Equal(t, enc, e.Encrypt(str), "何度暗号化しても同じ結果")
}

func TestDecrypt(t *testing.T) {
	iv := []byte{0x02, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01, 0x01}
	e, _ := New("AiZVGMsMrRGAFDeLhJKhLFSF7ULbHXhr", iv)
	equal1 := "あいうえお"
	equal2 := "abcde.123"
	str1 := e.Encrypt(equal1)
	str2 := e.Encrypt(equal2)
	assert.Equal(t, equal1, d(e, str1))
	assert.Equal(t, equal2, d(e, str2))
	assert.Equal(t, equal1, d(e, e.Encrypt(equal1)))
}

func d(e *Encrypt, v string) string {
	s, _ := e.Decrypt(v)
	return *s
}
