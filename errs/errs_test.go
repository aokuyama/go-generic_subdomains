package errs

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrors(t *testing.T) {
	errs := New()
	assert.NoError(t, errs.Err())
	errs.Append(errors.New("err"))
	assert.Error(t, errs.Err())
	errs.Append(errors.New("err"))
	assert.Error(t, errs.Err())
}

func TestAppendErrMessage(t *testing.T) {
	errs := New()
	errs.Append(errors.New("123"))
	errs.Append(errors.New("abc"))
	assert.Equal(t, "123\nabc\n", errs.Err().Error())
}

func TestAppendErrMessage2(t *testing.T) {
	var err error
	errs := New()
	err = errors.New("123")
	errs.Append(err)
	err = errors.New("abc")
	errs.Append(err)
	assert.Equal(t, "123\nabc\n", errs.Err().Error())
}

func TestAppendNil(t *testing.T) {
	errs := New()
	errs.Append(nil)
	assert.NoError(t, errs.Err())
	errs.Append(errors.New("abc"))
	assert.Error(t, errs.Err())
	assert.Equal(t, "abc\n", errs.Err().Error())
}
