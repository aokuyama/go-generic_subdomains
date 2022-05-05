package api

import (
	"encoding/json"
	"errors"
	"fmt"
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
	key, err := s.getDescribeKey()
	if err != nil {
		fmt.Println("failed start. response:", (string)(*s.StartApi.GetResult()))
		return err
	}
	fmt.Println("start succeeded. key:", *key)
	s.setFirstCallTime(s.Now())
	for {
		s.waitForNext(s.Now())
		err = s.doDescribeApi()
		if err != nil {
			return err
		}
		if !s.IsDescribeResponse() {
			fmt.Println("failed describe. response:", (string)(*s.DescribeApi.GetResult()))
			return errors.New("Failed describe.")
		}
		s.setCallTime(s.Now())
		if s.IsCompleted() {
			break
		}
	}
	if !s.IsSucceeded() {
		fmt.Println("failed api. response:", (string)(*s.DescribeApi.GetResult()))
		return errors.New("Failed Step API")
	}
	fmt.Println("describe succeeded.")
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

func (s *StepApi) GetStatusCode() int {
	return 0
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
	if len(jsonObj.ExecutionArn) == 0 {
		return nil, errors.New("キーの取得失敗")
	}
	key := jsonObj.ExecutionArn
	return &key, nil
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

func (s *StepApi) IsSucceeded() bool {
	if !s.DescribeApi.IsCompleted() {
		return false
	}
	var dr DescribeResponse
	json.Unmarshal(*s.DescribeApi.GetResult(), &dr)
	return dr.IsSucceeded()
}

func (s *StepApi) IsDescribeResponse() bool {
	if !s.DescribeApi.IsCompleted() {
		return false
	}
	var dr DescribeResponse
	json.Unmarshal(*s.DescribeApi.GetResult(), &dr)
	return dr.IsDescribeResponse()
}
