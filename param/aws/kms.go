package aws

import (
	"encoding/base64"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

type Kms struct {
	svc *kms.KMS
}

func NewKms(region string) *Kms {
	s := Kms{}
	se := session.Must(session.NewSession())
	s.svc = kms.New(se, aws.NewConfig().WithRegion(region))
	return &s
}

func (s *Kms) GetValue(v string) (string, error) {
	data, _ := base64.StdEncoding.DecodeString(v)

	input := &kms.DecryptInput{
		CiphertextBlob: []byte(data),
	}
	result, err := s.svc.Decrypt(input)
	if err != nil {
		return "", err
	}
	text, _ := base64.StdEncoding.DecodeString(base64.StdEncoding.EncodeToString(result.Plaintext))
	return string(text), nil
}

func (s *Kms) ConvertValue(value string) (string, error) {
	key := s.parseEncrypted(value)
	if len(key) == 0 {
		return value, nil
	}
	return s.GetValue(key)
}

func (s *Kms) parseEncrypted(v string) string {
	if strings.HasPrefix(v, "#KMS#") {
		r := string([]rune(v)[5:])
		return r
	}
	return ""
}
