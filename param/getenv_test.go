package param

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotConvertNormalValue(t *testing.T) {
	assert.Equal(t, "abc", ConvertValue("abc"))
	assert.Equal(t, "_1234", ConvertValue("_1234"))
	assert.Equal(t, "#あいうえお", ConvertValue("#あいうえお"))
}
