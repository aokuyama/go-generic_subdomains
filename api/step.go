package api

import (
	"encoding/json"
	"errors"
	"time"
)

type Step struct {
	IntervalSecFirst int
	IntervalSec      int
	next_call        *time.Time
	count_wait       int
	result           *[]byte
	StartApi         *Api
	DescribeApi      *Api
}

func NewStep(start_url string, describe_url string) *Step {
	start, _ := NewApi(start_url)
	describe, _ := NewApi(describe_url)

	s := Step{
		IntervalSecFirst: 30,
		IntervalSec:      10,
		next_call:        nil,
		count_wait:       0,
		StartApi:         start,
		DescribeApi:      describe,
	}
	return &s
}

func (s *Step) Do(body interface{}) error {
	var err error
	err = s.doStartApi(body)
	if err != nil {
		return err
	}
	s.setFirstCallTime(s.Now())
	for {
		s.waitForNext(s.Now())
		err = s.doDescribeApi()
		if err != nil {
			return err
		}
		s.setCallTime(s.Now())
		if s.isStepApiCompleted() {
			return nil
		}
	}
	return nil
}

func (s *Step) GetResult() *[]byte {
	return s.DescribeApi.GetResult()
}

func (s *Step) GetJson() interface{} {
	var jsonObj DescribeResponse
	_ = json.Unmarshal(*s.GetResult(), &jsonObj)
	return jsonObj.Output
}

func (s *Step) doStartApi(body interface{}) error {
	if s.StartApi.isCompleted() {
		return errors.New("startの二重実行")
	}
	return s.StartApi.Do(body)
}

func (s *Step) getDescribeKey() (*string, error) {
	if !s.StartApi.isCompleted() {
		return nil, errors.New("startの実行前")
	}
	var jsonObj StartApiResponse
	err := json.Unmarshal(*s.StartApi.GetResult(), &jsonObj)
	if err != nil {
		return nil, err
	}
	return &jsonObj.ExecutionArn, nil
}

func (s *Step) doDescribeApi() error {
	key, err := s.getDescribeKey()
	if err != nil {
		return err
	}
	dr := NewDescribeRequest(*key)
	return s.DescribeApi.Do(dr)
}

func (s *Step) isStepApiCompleted() bool {
	if !s.DescribeApi.isCompleted() {
		return false
	}
	var dr DescribeResponse
	json.Unmarshal(*s.DescribeApi.GetResult(), &dr)
	return dr.isCompleted()
}
