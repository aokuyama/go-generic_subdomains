package event

import "github.com/aokuyama/go-generic_subdomains/errs"

type Publisher struct {
	subsclibers []Subscriber
}

func NewPublisher() *Publisher {
	p := Publisher{}
	return &p
}

func (p *Publisher) Regist(s ...Subscriber) {
	p.subsclibers = append(p.subsclibers, s...)
}

func (p *Publisher) Publish(e Event) error {
	errs := errs.New()
	for _, s := range p.subsclibers {
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
