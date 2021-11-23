package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type Store struct {
	ssm *ssm.SSM
}

func NewStore(region string) *Store {
	c := Store{}
	se := session.Must(session.NewSession())
	c.ssm = ssm.New(
		se,
		aws.NewConfig().WithRegion(region),
	)
	return &c
}

func (s *Store) GetValue(key string) (*string, error) {
	res, err := s.ssm.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		return nil, err
	}
	value := res.Parameter.Value
	return value, nil
}
