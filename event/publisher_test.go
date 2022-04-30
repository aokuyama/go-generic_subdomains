package event_test

import (
	"testing"

	. "github.com/aokuyama/go-generic_subdomains/event"
	"github.com/stretchr/testify/assert"
)

func TestPublish(t *testing.T) {
	suba := SubscriberA{}
	subb := SubscriberB{}
	subw := SubscriberW{}
	pub := NewPublisher()
	pub.Register(&suba)
	pub.Register(&subb, &subw)
	e1 := ExecOrderEvent{}
	assert.Equal(t, 0, suba.Count())
	assert.Equal(t, 0, subb.Count())
	assert.Equal(t, 0, subw.Count())
	pub.Publish(&e1)
	assert.Equal(t, 1, suba.Count())
	assert.Equal(t, 0, subb.Count())
	assert.Equal(t, 1, subw.Count())
	e2 := GetResultEvent{}
	pub.Publish(&e2)
	assert.Equal(t, 1, suba.Count())
	assert.Equal(t, 1, subb.Count())
	assert.Equal(t, 2, subw.Count())
}
