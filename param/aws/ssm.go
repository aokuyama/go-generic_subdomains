package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type Ssm struct {
	ssm *ssm.SSM
}

func NewSsm(region string) *Ssm {
	c := Ssm{}
	se := session.Must(session.NewSession())
	c.ssm = ssm.New(
		se,
		aws.NewConfig().WithRegion(region),
	)
	return &c
}

func (s *Ssm) GetValue(key string) (string, error) {
	res, err := s.ssm.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", err
	}
	value := res.Parameter.Value
	return *value, nil
}

func (s *Ssm) ConvertValue(value string) (string, error) {
	key := s.parseStoreKey(value)
	if len(key) == 0 {
		return value, nil
	}
	return s.GetValue(key)
}

func (s *Ssm) parseStoreKey(v string) string {
	if strings.HasPrefix(v, "#SSM#") {
		r := string([]rune(v)[5:])
		return r
	}
	return ""
}
