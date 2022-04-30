//go:generate mockgen -source=$GOFILE -package=mock -destination=mock/$GOFILE
package event

import (
	"reflect"
	"regexp"
)

type Event interface {
	Data() string
}

func IsSubscribe(s Subscriber, e Event) bool {
	st := s.Type()
	subscribe_type := reflect.TypeOf(st)
	event_type := reflect.TypeOf(e)
	if subscribe_type.String() == "string" {
		v := reflect.ValueOf(st)
		if v.String() == "" {
			return false
		}
		r := regexp.MustCompile(v.String())
		return r.MatchString(event_type.String())
	}
	return subscribe_type == event_type
}
