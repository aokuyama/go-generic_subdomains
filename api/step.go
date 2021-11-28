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
	StartApi         Api
	DescribeApi      Api
}

func NewStepApi(start Api, describe Api) *StepApi {
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
	if !s.StartApi.IsCompleted() {
		err = s.doStartApi(body)
	}
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
		if s.IsCompleted() {
			return nil
		}
	}
	return nil
}

func (s *StepApi) GetResult() *[]byte {
	var dr DescribeResponse
	err := json.Unmarshal(*s.DescribeApi.GetResult(), &dr)
	if err != nil {
		panic(err)
	}
	b := []byte(dr.Output)
	return &b
}

func (s *StepApi) doStartApi(body interface{}) error {
	if s.StartApi.IsCompleted() {
		return errors.New("startの二重実行")
	}
	return s.StartApi.Do(body)
}

func (s *StepApi) getDescribeKey() (*string, error) {
	if !s.StartApi.IsCompleted() {
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

func (s *StepApi) IsCompleted() bool {
	if !s.DescribeApi.IsCompleted() {
		return false
	}
	var dr DescribeResponse
	json.Unmarshal(*s.DescribeApi.GetResult(), &dr)
	return dr.IsCompleted()
}
