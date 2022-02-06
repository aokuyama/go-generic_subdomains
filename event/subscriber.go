package event

type Subscriber interface {
	Type() interface{}
	Subscribe(e Event) error
}
