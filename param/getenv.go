package param

import (
	"os"

	"github.com/aokuyama/go-generic_subdomains/param/aws"
)

type Store interface {
	GetValue(key string) (*string, error)
	GetValueIfKey(value string) (*string, error)
}

var ssm Store
var kms Store

func Getenv(s string) string {
	env := os.Getenv(s)
	if ssm == nil {
		ssm = aws.NewSsm("ap-northeast-1")
	}
	value, err := ssm.GetValueIfKey(env)
	if err != nil {
		panic(err)
	}
	if kms == nil {
		kms = aws.NewKms("ap-northeast-1")
	}
	value2, err := kms.GetValueIfKey(*value)
	if err != nil {
		panic(err)
	}
	return *value2
}
