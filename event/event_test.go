package event_test

import (
	"testing"

	. "github.com/aokuyama/go-generic_subdomains/event"
	"github.com/stretchr/testify/assert"
)

type ExecOrderEvent struct {
}

func (e *ExecOrderEvent) Data() string {
	return ""
}

type GetResultEvent struct {
}

func (e *GetResultEvent) Data() string {
	return ""
}

type SubscriberA struct {
	count int
}

func (s *SubscriberA) Type() interface{} {
	return &ExecOrderEvent{}
}
func (s *SubscriberA) Subscribe(e Event) error {
	s.count += 1
	return nil
}
func (s *SubscriberA) Count() int {
	return s.count
}

type SubscriberB struct {
	count int
}

func (s *SubscriberB) Type() interface{} {
	return &GetResultEvent{}
}
func (s *SubscriberB) Subscribe(e Event) error {
	s.count += 1
	return nil
}
func (s *SubscriberB) Count() int {
	return s.count
}

type SubscriberW struct {
	count int
}

func (s *SubscriberW) Type() interface{} {
	return ".*"
}
func (s *SubscriberW) Subscribe(e Event) error {
	s.count += 1
	return nil
}
func (s *SubscriberW) Count() int {
	return s.count
}

type SubscriberN struct {
	count int
}

func (s *SubscriberN) Type() interface{} {
	return ""
}
func (s *SubscriberN) Subscribe(e Event) error {
	return nil
}
func (s *SubscriberN) Count() int {
	return s.count
}

func TestIsSubscribe(t *testing.T) {
	assert.True(t, IsSubscribe(&SubscriberA{}, &ExecOrderEvent{}))
	assert.False(t, IsSubscribe(&SubscriberA{}, &GetResultEvent{}))
	assert.False(t, IsSubscribe(&SubscriberB{}, &ExecOrderEvent{}))
	assert.True(t, IsSubscribe(&SubscriberB{}, &GetResultEvent{}))
	assert.True(t, IsSubscribe(&SubscriberW{}, &ExecOrderEvent{}))
	assert.True(t, IsSubscribe(&SubscriberW{}, &GetResultEvent{}))
	assert.False(t, IsSubscribe(&SubscriberN{}, &ExecOrderEvent{}))
	assert.False(t, IsSubscribe(&SubscriberN{}, &GetResultEvent{}))
}
