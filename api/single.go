package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type SingleApi struct {
	TimeOut     int
	Method      string
	ContentType string
	XApiKey     string
	url         url.URL
	result      *[]byte
}

func NewSingleApi(u string) (*SingleApi, error) {
	up, err := url.Parse(u)
	if err != nil {
		return nil, err
	}
	s := SingleApi{
		TimeOut:     10,
		Method:      "POST",
		ContentType: "application/json",
		url:         *up,
	}
	return &s, nil
}

func (a *SingleApi) Do(body interface{}) error {
	client := &http.Client{}
	client.Timeout = time.Second * time.Duration(a.TimeOut)
	var err error

	req, err := a.createRequest(body)
	if err != nil {
		return err
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	r, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	a.result = &r
	return nil
}

func (a *SingleApi) createRequest(body interface{}) (*http.Request, error) {
	var b io.Reader
	if body == nil {
		b = nil
	} else {
		j, _ := json.Marshal(body)
		b = bytes.NewBuffer(j)
	}
	req, err := http.NewRequest(a.Method, a.url.String(), b)
	header := http.Header{}
	header.Add("Content-Type", a.ContentType)
	if len(a.XApiKey) > 0 {
		header.Add("x-api-key", a.XApiKey)
	}
	req.Header = header
	return req, err
}

func (a *SingleApi) GetResult() *[]byte {
	return a.result
}

func (a *SingleApi) isCompleted() bool {
	return a.GetResult() != nil
}
