package param

import (
	"os"

	"github.com/aokuyama/go-generic_subdomains/param/aws"
)

type Store interface {
	GetValue(key string) (string, error)
	ConvertValue(value string) (string, error)
}

var ssm Store
var kms Store

func Getenv(s string) string {
	return ConvertValue(os.Getenv(s))
}

func ConvertValue(s string) string {
	if ssm == nil {
		ssm = aws.NewSsm("ap-northeast-1")
	}
	s, err := ssm.ConvertValue(s)
	if err != nil {
		panic(err)
	}
	if kms == nil {
		kms = aws.NewKms("ap-northeast-1")
	}
	s, err = kms.ConvertValue(s)
	if err != nil {
		panic(err)
	}
	return s
}
