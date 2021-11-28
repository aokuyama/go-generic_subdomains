package api

import (
	"fmt"
	"time"
)

func (s *StepApi) Now() *time.Time {
	n := time.Now()
	return &n
}

func (s *StepApi) GetNextCall() (t *time.Time) {
	return s.next_call
}

func (s *StepApi) setFirstCallTime(t *time.Time) {
	t2 := t.Add(time.Second * time.Duration(s.IntervalSecFirst))
	s.next_call = &t2
}

func (s *StepApi) setCallTime(t *time.Time) {
	t2 := t.Add(time.Second * time.Duration(s.IntervalSec))
	s.next_call = &t2
}

func (s *StepApi) isCallableTime(now *time.Time) bool {
	return now.Equal(*s.next_call) || now.After(*s.next_call)
}

func (s *StepApi) GetInterval(now *time.Time) time.Duration {
	sub := s.next_call.Sub(*now)
	if sub < 0 {
		return time.Duration(0)
	}
	return sub
}

func (s *StepApi) waitForNext(now *time.Time) {
	w := s.GetInterval(now)
	fmt.Println("wait", w, "...")
	time.Sleep(w)
}
