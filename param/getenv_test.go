package param

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_storeKey(t *testing.T) {
	assert.Equal(t, "abc", *parseStoreKey("#SSM#abc"))
	assert.Equal(t, "_1234", *parseStoreKey("#SSM#_1234"))
	assert.Equal(t, "あいうえお", *parseStoreKey("#SSM#あいうえお"))
}

func Test_notStoreKey(t *testing.T) {
	var equal *string
	assert.Equal(t, equal, parseStoreKey("http://example.com"))
	assert.Equal(t, equal, parseStoreKey("#ssm#abc"))
	assert.Equal(t, equal, parseStoreKey("##SSM#abc"))
	assert.Equal(t, equal, parseStoreKey("SSM#abc"))
}
