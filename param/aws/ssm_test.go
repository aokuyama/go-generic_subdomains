package aws

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreKey(t *testing.T) {
	s := Ssm{}
	assert.Equal(t, "abc", *s.parseStoreKey("#SSM#abc"))
	assert.Equal(t, "_1234", *s.parseStoreKey("#SSM#_1234"))
	assert.Equal(t, "あいうえお", *s.parseStoreKey("#SSM#あいうえお"))
}

func TestNotStoreKey(t *testing.T) {
	s := Ssm{}
	var equal *string
	assert.Equal(t, equal, s.parseStoreKey("http://example.com"))
	assert.Equal(t, equal, s.parseStoreKey("#ssm#abc"))
	assert.Equal(t, equal, s.parseStoreKey("##SSM#abc"))
	assert.Equal(t, equal, s.parseStoreKey("SSM#abc"))
	assert.Equal(t, equal, s.parseStoreKey("#KMS#abc"))
}
