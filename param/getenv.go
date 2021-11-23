package param

import (
	"github.com/aokuyama/go-generic_subdomains/param/aws"
	"os"
	"strings"
)

type Store interface {
	GetValue(key string) (*string, error)
}

var store Store

func Getenv(s string) string {
	env := os.Getenv(s)
	if store == nil {
		store = aws.NewStore("ap-northeast-1")
	}
	return replaceValue(env, store)
}

func replaceValue(s string, store Store) string {
	key := parseStoreKey(s)
	if key == nil {
		return s
	}
	value, err := store.GetValue(*key)
	if err != nil {
		panic(err)
	}
	return *value
}

func parseStoreKey(s string) *string {
	if strings.HasPrefix(s, "#SSM#") {
		r := string([]rune(s)[5:])
		return &r
	}
	return nil
}
