package api

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type AwsLambda struct {
	svc      *lambda.Lambda
	function string
	result   *[]byte
}

func New(function string) *AwsLambda {
	s := session.New()
	svc := lambda.New(s, aws.NewConfig())
	return &AwsLambda{
		svc:      svc,
		function: function,
	}
}

func (l *AwsLambda) Do(body interface{}) error {
	jsonBytes, _ := json.Marshal(body)

	input := &lambda.InvokeInput{
		FunctionName: aws.String(l.function),
		Payload:      jsonBytes,
	}
	r, err := l.svc.Invoke(input)

	if err != nil {
		return err
	}
	re := []byte(r.String())
	l.result = &re
	return nil
}

func (l *AwsLambda) GetResult() *[]byte {
	return l.result
}

func (l *AwsLambda) GetStatusCode() int {
	return 0
}

func (l *AwsLambda) IsCompleted() bool {
	return l.result != nil
}
