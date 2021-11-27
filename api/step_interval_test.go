package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_setFirstCallTime(t *testing.T) {
	var tm time.Time
	var equal time.Time
	s := initStep()
	s.IntervalSecFirst = 30
	tm = time.Date(2021, 11, 8, 14, 24, 0, 0, time.UTC)
	s.setFirstCallTime(&tm)
	equal = time.Date(2021, 11, 8, 14, 24, 30, 0, time.UTC)
	assert.Equal(t, &equal, s.GetNextCall(), "次の問い合わせ時間")
	s.IntervalSecFirst = 15
	tm = time.Date(2021, 11, 8, 14, 25, 10, 0, time.UTC)
	s.setFirstCallTime(&tm)
	equal = time.Date(2021, 11, 8, 14, 25, 25, 0, time.UTC)
	assert.Equal(t, &equal, s.GetNextCall(), "次の問い合わせ時間")
}

func Test_setCallTime(t *testing.T) {
	var tm time.Time
	var equal time.Time
	s := initStep()
	s.IntervalSec = 5
	tm = time.Date(2021, 11, 8, 14, 24, 0, 0, time.UTC)
	s.setCallTime(&tm)
	equal = time.Date(2021, 11, 8, 14, 24, 5, 0, time.UTC)
	assert.Equal(t, &equal, s.GetNextCall(), "次の問い合わせ時間")
	s.IntervalSec = 20
	tm = time.Date(2021, 11, 8, 14, 25, 10, 0, time.UTC)
	s.setCallTime(&tm)
	equal = time.Date(2021, 11, 8, 14, 25, 30, 0, time.UTC)
	assert.Equal(t, &equal, s.GetNextCall(), "次の問い合わせ時間")
}

func Test_isCallableTime(t *testing.T) {
	var tm time.Time
	var now time.Time
	s := initStep()
	s.IntervalSecFirst = 5
	tm = time.Date(2021, 11, 8, 14, 24, 0, 0, time.UTC)
	s.setFirstCallTime(&tm)
	now = time.Date(2021, 11, 8, 14, 24, 4, 0, time.UTC)
	assert.Equal(t, false, s.isCallableTime(&now), "問い合わせ時間前")
	now = time.Date(2021, 11, 8, 14, 24, 5, 0, time.UTC)
	assert.Equal(t, true, s.isCallableTime(&now), "問い合わせ時間ちょうど")
	now = time.Date(2021, 11, 8, 14, 24, 6, 0, time.UTC)
	assert.Equal(t, true, s.isCallableTime(&now), "問い合わせ時間以降")
}

func Test_getNextCallTimeInterval(t *testing.T) {
	var now time.Time
	var tm time.Time
	s := initStep()
	s.IntervalSecFirst = 0
	s.IntervalSec = 0
	tm = time.Date(2021, 11, 8, 14, 24, 0, 0, time.UTC)
	s.setFirstCallTime(&tm)
	now = time.Date(2021, 11, 8, 14, 23, 0, 0, time.UTC)
	assert.Equal(t, (60 * time.Second), s.GetInterval(&now), "1分前")
	now = time.Date(2021, 11, 8, 14, 23, 50, 0, time.UTC)
	assert.Equal(t, (10 * time.Second), s.GetInterval(&now), "10秒前")
	now = time.Date(2021, 11, 8, 14, 24, 0, 0, time.UTC)
	assert.Equal(t, (0 * time.Second), s.GetInterval(&now), "ちょうど")
	now = time.Date(2021, 11, 8, 14, 26, 0, 0, time.UTC)
	assert.Equal(t, (0 * time.Second), s.GetInterval(&now), "時間が過ぎた")
}
