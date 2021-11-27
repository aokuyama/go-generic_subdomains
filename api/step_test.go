package api

import (
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"testing"
	//"fmt"
)

func initStep() *Step {
	s := NewStep(
		"http://start.mock.example.com",
		"http://describe.mock.example.com",
	)
	return s
}

func setStarted(s *Step) {
	result := []byte("{\"executionArn\":\"arn:aws:states:ap-northeast-1:xxxx:xxxx:xxxx:1234\",\"startDate\":1.627105289676E9}")
	s.StartApi.result = &result
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
	assert.Equal(t, false, s.isStepApiCompleted(), "実行前")
	setStarted(s)
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("POST", "http://describe.mock.example.com",
		httpmock.NewStringResponder(200, "mocked describe api"),
	)
	s.doDescribeApi()
	assert.Equal(t, false, s.isStepApiCompleted(), "statusが含まれない場合")
	httpmock.RegisterResponder("POST", "http://describe.mock.example.com",
		httpmock.NewStringResponder(200, "{\"status\":\"SUCCEEDED\"}"),
	)
	s.doDescribeApi()
	assert.Equal(t, true, s.isStepApiCompleted(), "成功した場合")
	httpmock.RegisterResponder("POST", "http://describe.mock.example.com",
		httpmock.NewStringResponder(200, "{\"status\":\"PENDING\"}"),
	)
	s.doDescribeApi()
	assert.Equal(t, false, s.isStepApiCompleted(), "継続中の場合")
	httpmock.RegisterResponder("POST", "http://describe.mock.example.com",
		httpmock.NewStringResponder(200, "{\"status\":\"FAILED\"}"),
	)
	s.doDescribeApi()
	assert.Equal(t, true, s.isStepApiCompleted(), "失敗した場合")
}
