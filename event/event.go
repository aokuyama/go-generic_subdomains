package event

import (
	"reflect"
	"regexp"
)

type Event interface {
	Data() string
}

func IsSubscribe(s Subscriber, e Event) bool {
	subscribe_type := reflect.TypeOf(s.Type())
	event_type := reflect.TypeOf(e)
	if subscribe_type.String() == "string" {
		v := reflect.ValueOf(s.Type())
		r := regexp.MustCompile(v.String())
		return r.MatchString(event_type.String())
	}
	return subscribe_type == event_type
}
