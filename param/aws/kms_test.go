package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptedValue(t *testing.T) {
	s := Kms{}
	assert.Equal(t, "abz", s.parseEncrypted("#KMS#abz"))
	assert.Equal(t, "_124", s.parseEncrypted("#KMS#_124"))
}

func TestNotEncryptedValue(t *testing.T) {
	s := Kms{}
	assert.Equal(t, "", s.parseEncrypted("#SSM#abz"))
	assert.Equal(t, "", s.parseEncrypted("KMS#_124"))
}
