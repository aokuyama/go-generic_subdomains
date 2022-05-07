package queue

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type AwsSqs struct {
	client *sqs.SQS
	url    string
}

func New(url string) (*AwsSqs, error) {
	se, err := session.NewSession()
	if err != nil {
		return nil, err
	}
	return &AwsSqs{client: sqs.New(se), url: url}, nil
}

func (s *AwsSqs) Publish(body string, group_id string) (*string, error) {
	smi := sqs.SendMessageInput{
		QueueUrl:       &s.url,
		MessageBody:    &body,
		MessageGroupId: &group_id,
	}
	smo, err := s.client.SendMessage(&smi)
	if err != nil {
		return nil, err
	}
	r := smo.String()
	return &r, nil
}
