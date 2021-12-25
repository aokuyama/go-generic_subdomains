package database

import (
	"reflect"
	"strings"
	"unicode"

	"github.com/gertd/go-pluralize"
)

func TableName(st interface{}) string {
	return toPlural(toSnakeCase(structName(st)))
}

func structName(st interface{}) string {
	if t := reflect.TypeOf(st); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func toSnakeCase(s string) string {
	b := &strings.Builder{}
	for i, r := range s {
		if i == 0 {
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		if unicode.IsUpper(r) {
			b.WriteRune('_')
			b.WriteRune(unicode.ToLower(r))
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func toPlural(s string) string {
	plu := pluralize.NewClient()
	return plu.Plural(s)
}
