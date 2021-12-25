package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructToTableName(t *testing.T) {
	type TestStruct struct{}
	type Test1Order struct{}
	assert.Equal(t, "test_structs", TableName(TestStruct{}), "スネークケースかつ複数形")
	assert.Equal(t, "test1_orders", TableName(Test1Order{}), "スネークケースかつ複数形")
}
