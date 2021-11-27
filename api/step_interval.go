package api

import (
	"fmt"
	"time"
)

func (s *Step) Now() *time.Time {
	n := time.Now()
	return &n
}

func (s *Step) GetNextCall() (t *time.Time) {
	return s.next_call
}

func (s *Step) setFirstCallTime(t *time.Time) {
	t2 := t.Add(time.Second * time.Duration(s.IntervalSecFirst))
	s.next_call = &t2
}

func (s *Step) setCallTime(t *time.Time) {
	t2 := t.Add(time.Second * time.Duration(s.IntervalSec))
	s.next_call = &t2
}

func (s *Step) isCallableTime(now *time.Time) bool {
	return now.Equal(*s.next_call) || now.After(*s.next_call)
}

func (s *Step) GetInterval(now *time.Time) time.Duration {
	sub := s.next_call.Sub(*now)
	if sub < 0 {
		return time.Duration(0)
	}
	return sub
}

func (s *Step) waitForNext(now *time.Time) {
	w := s.GetInterval(now)
	fmt.Println("wait", w, "...")
	time.Sleep(w)
}
