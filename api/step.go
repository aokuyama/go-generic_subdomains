package api

import (
	"encoding/json"
	"errors"
	"time"
)

type StepApi struct {
	IntervalSecFirst int
	IntervalSec      int
	next_call        *time.Time
	count_wait       int
	result           *[]byte
	StartApi         *SingleApi
	DescribeApi      *SingleApi
}

func NewStepApi(start_url string, describe_url string) *StepApi {
	start, _ := NewSingleApi(start_url)
	describe, _ := NewSingleApi(describe_url)

	s := StepApi{
		IntervalSecFirst: 30,
		IntervalSec:      10,
		next_call:        nil,
		count_wait:       0,
		StartApi:         start,
		DescribeApi:      describe,
	}
	return &s
}

func (s *StepApi) Do(body interface{}) error {
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

func (s *StepApi) GetResult() *[]byte {
	return s.DescribeApi.GetResult()
}

func (s *StepApi) GetJson() interface{} {
	var jsonObj DescribeResponse
	_ = json.Unmarshal(*s.GetResult(), &jsonObj)
	return jsonObj.Output
}

func (s *StepApi) doStartApi(body interface{}) error {
	if s.StartApi.isCompleted() {
		return errors.New("startの二重実行")
	}
	return s.StartApi.Do(body)
}

func (s *StepApi) getDescribeKey() (*string, error) {
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

func (s *StepApi) doDescribeApi() error {
	key, err := s.getDescribeKey()
	if err != nil {
		return err
	}
	dr := NewDescribeRequest(*key)
	return s.DescribeApi.Do(dr)
}

func (s *StepApi) isStepApiCompleted() bool {
	if !s.DescribeApi.isCompleted() {
		return false
	}
	var dr DescribeResponse
	json.Unmarshal(*s.DescribeApi.GetResult(), &dr)
	return dr.isCompleted()
}
