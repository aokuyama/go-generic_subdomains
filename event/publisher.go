package event

import "github.com/aokuyama/go-generic_subdomains/errs"

type Publisher struct {
	subscribers []Subscriber
}

func NewPublisher() *Publisher {
	p := Publisher{}
	return &p
}

func (p *Publisher) Register(s ...Subscriber) {
	p.subscribers = append(p.subscribers, s...)
}

func (p *Publisher) Publish(e Event) error {
	errs := errs.New()
	for _, s := range p.subscribers {
		errs.Append(p.delivery(s, e))
	}
	return errs.Err()
}

func (p *Publisher) delivery(s Subscriber, e Event) error {
	if !IsSubscribe(s, e) {
		return nil
	}
	return s.Subscribe(e)
}
