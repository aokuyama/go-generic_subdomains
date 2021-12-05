package api

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func initStep() *StepApi {
	start, _ := NewSingleApi("http://start.mock.example.com")
	describe, _ := NewSingleApi("http://describe.mock.example.com")
	s := NewStepApi(
		start,
		describe,
	)
	s.IntervalSec = 0
	s.IntervalSecFirst = 0
	return s
}

func setStarted(s *StepApi) {
	start, _ := NewSingleApi("http://start.mock.example.com")
	result := []byte("{\"executionArn\":\"arn:aws:states:ap-northeast-1:xxxx:xxxx:xxxx:1234\",\"startDate\":1.627105289676E9}")
	start.result = &result
	s.StartApi = start
}

func Test_startApi(t *testing.T) {
	s := initStep()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://start.mock.example.com",
		httpmock.NewStringResponder(200, "mocked start api"),
	)
	s.doStartApi(nil)
	assert.Equal(t, "mocked start api", string(*s.StartApi.GetResult()), "startApiの結果を受け取る")
	err := s.doStartApi(nil)
	assert.Error(t, err, "二重実行の禁止")
}

func Test_getDescribeKey(t *testing.T) {
	s := initStep()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://start.mock.example.com",
		httpmock.NewStringResponder(200, "{\"executionArn\":\"arn:aws:states:ap-northeast-1:xxxx:xxxx:xxxx:1234\",\"startDate\":1.627105289676E9}"),
	)
	_, err := s.getDescribeKey()
	assert.Error(t, err, "startの実行前は取得できない")
	s.doStartApi(nil)
	key, _ := s.getDescribeKey()
	assert.Equal(t, "arn:aws:states:ap-northeast-1:xxxx:xxxx:xxxx:1234", string(*key), "startApiの結果からexcuteArnを取得")
}

func Test_FailgetDescribeKey(t *testing.T) {
	s := initStep()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://start.mock.example.com",
		httpmock.NewStringResponder(200, "NG"),
	)
	s.doStartApi(nil)
	_, err := s.getDescribeKey()
	assert.Error(t, err, "startの実行結果が型に合わない場合エラー")
}
func Test_FailgetDescribeKey2(t *testing.T) {
	s := initStep()
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://start.mock.example.com",
		httpmock.NewStringResponder(200, `{"message":"Forbidden"}`),
	)
	s.doStartApi(nil)
	_, err := s.getDescribeKey()
	assert.Error(t, err, "startの実行結果、キーが取得できていないとエラー")
}

func Test_DoDescribeApi(t *testing.T) {
	s := initStep()
	err := s.doDescribeApi()
	assert.Error(t, err, "startApiの実行前はエラー")
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://start.mock.example.com",
		httpmock.NewStringResponder(200, "{\"executionArn\":\"arn:aws:states:ap-northeast-1:xxxx:xxxx:xxxx:1234\",\"startDate\":1.627105289676E9}"),
	)
	httpmock.RegisterResponder("POST", "http://describe.mock.example.com",
		httpmock.NewStringResponder(200, "mocked describe api"),
	)
	s.doStartApi(nil)
	s.doDescribeApi()
	assert.Equal(t, "mocked describe api", string(*s.DescribeApi.GetResult()), "describeApiの結果を受け取る")
}

func Test_DoDescribeResult(t *testing.T) {
	s := initStep()
	assert.Equal(t, false, s.IsCompleted(), "実行前")
	setStarted(s)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://describe.mock.example.com",
		httpmock.NewStringResponder(200, "mocked describe api"),
	)
	s.doDescribeApi()
	assert.Equal(t, false, s.IsCompleted(), "statusが含まれない場合")
	httpmock.RegisterResponder("POST", "http://describe.mock.example.com",
		httpmock.NewStringResponder(200, "{\"status\":\"SUCCEEDED\"}"),
	)
	s.doDescribeApi()
	assert.Equal(t, true, s.IsCompleted(), "成功した場合")
	httpmock.RegisterResponder("POST", "http://describe.mock.example.com",
		httpmock.NewStringResponder(200, "{\"status\":\"PENDING\"}"),
	)
	s.doDescribeApi()
	assert.Equal(t, false, s.IsCompleted(), "継続中の場合")
	httpmock.RegisterResponder("POST", "http://describe.mock.example.com",
		httpmock.NewStringResponder(200, "{\"status\":\"FAILED\"}"),
	)
	s.doDescribeApi()
	assert.Equal(t, true, s.IsCompleted(), "失敗した場合")
}

func TestGetResult(t *testing.T) {
	s := initStep()
	setStarted(s)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://start.mock.example.com",
		httpmock.NewStringResponder(200, "{\"executionArn\":\"arn:aws:states:ap-northeast-1:xxxx:xxxx:xxxx:1234\",\"startDate\":1.627105289676E9}"),
	)
	result := `{"executionArn":"arn:aws:states:ap-northeast-1:xxxx:xxxx:xxxx:1234","input":"{\"param\":\"a\"}","output":"{\"error\": 0, \"list\": [{\"date\": \"2021-12-29\", \"value\": \"1234\"}, {\"date\": \"2021-12-27\", \"value\": \"abcd\"}]}","status":"SUCCEEDED"}`
	httpmock.RegisterResponder("POST", "http://describe.mock.example.com",
		httpmock.NewStringResponder(200, result),
	)
	err := s.Do(nil)
	assert.NoError(t, err, "成功")
	equal := `{"error": 0, "list": [{"date": "2021-12-29", "value": "1234"}, {"date": "2021-12-27", "value": "abcd"}]}`
	assert.Equal(t, equal, string(*s.GetResult()), "statusが含まれない場合")
}

func TestDoFailed(t *testing.T) {
	s := initStep()
	setStarted(s)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://start.mock.example.com",
		httpmock.NewStringResponder(200, "{\"executionArn\":\"arn:aws:states:ap-northeast-1:xxxx:xxxx:xxxx:1234\",\"startDate\":1.627105289676E9}"),
	)
	result := `{"executionArn":"arn:aws:states:ap-northeast-1:xxxx:xxxx:xxxx:1234","input":"{\"param\":\"a\"}","output":"{\"error\": 0, \"list\": [{\"date\": \"2021-12-29\", \"value\": \"1234\"}, {\"date\": \"2021-12-27\", \"value\": \"abcd\"}]}","status":"FAILED"}`
	httpmock.RegisterResponder("POST", "http://describe.mock.example.com",
		httpmock.NewStringResponder(200, result),
	)
	err := s.Do(nil)
	assert.Error(t, err, "APIの実行結果がFAIL")
}

func TestDescribeErrorJsonStatus(t *testing.T) {
	s := initStep()
	setStarted(s)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://start.mock.example.com",
		httpmock.NewStringResponder(200, "{\"executionArn\":\"arn:aws:states:ap-northeast-1:xxxx:xxxx:xxxx:1234\",\"startDate\":1.627105289676E9}"),
	)
	result := `{"message":"Forbidden"}`
	httpmock.RegisterResponder("POST", "http://describe.mock.example.com",
		httpmock.NewStringResponder(200, result),
	)
	err := s.Do(nil)
	assert.Error(t, err, "Describeの実行結果、statusが取得できない")
}

func TestDescribeErrorStatus(t *testing.T) {
	s := initStep()
	setStarted(s)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://start.mock.example.com",
		httpmock.NewStringResponder(200, "{\"executionArn\":\"arn:aws:states:ap-northeast-1:xxxx:xxxx:xxxx:1234\",\"startDate\":1.627105289676E9}"),
	)
	result := "NG"
	httpmock.RegisterResponder("POST", "http://describe.mock.example.com",
		httpmock.NewStringResponder(200, result),
	)
	err := s.Do(nil)
	assert.Error(t, err, "Describeの実行結果、statusが取得できない")
}
