package api

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"net/http/httputil"
	"testing"
)

func Test_Initialize(t *testing.T) {
	a, _ := NewApi("http://example.com")
	assert.Equal(t, true, a.GetResult() == nil)
}

func Test_CreateRequest(t *testing.T) {
	a, _ := NewApi("http://example.com")
	req, _ := a.createRequest(nil)
	dump, _ := httputil.DumpRequestOut(req, true)
	e := "POST / HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 0\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n"
	assert.Equal(t, e, string(dump))
}

type TestRequestBody struct {
	Id   int    `json: "id"`
	Name string `json: "name"`
}

func Test_CreateJsonPostRequest(t *testing.T) {
	a, _ := NewApi("http://example.com")
	s := TestRequestBody{
		Id:   1,
		Name: "aaaa",
	}
	req, _ := a.createRequest(&s)
	dump, _ := httputil.DumpRequestOut(req, true)
	e := "POST / HTTP/1.1\r\nHost: example.com\r\nUser-Agent: Go-http-client/1.1\r\nContent-Length: 22\r\nContent-Type: application/json\r\nAccept-Encoding: gzip\r\n\r\n{\"Id\":1,\"Name\":\"aaaa\"}"
	assert.Equal(t, e, string(dump))
}

func Test_GetResult(t *testing.T) {
	a, _ := NewApi("http://mock.example.com")
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://mock.example.com",
		httpmock.NewStringResponder(200, "mocked"),
	)
	assert.Equal(t, false, a.isCompleted())
	a.Do(nil)
	assert.Equal(t, "mocked", string(*a.GetResult()))
	assert.Equal(t, true, a.isCompleted())
}
