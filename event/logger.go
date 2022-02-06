package event

import (
	"fmt"
	"reflect"
)

type Logger struct {
}

func (s *Logger) Type() interface{} {
	return ".*"
}

func (s *Logger) Subscribe(e Event) error {
	fmt.Println(reflect.TypeOf(e), e.Data())
	return nil
}
