//go:generate mockgen -source=$GOFILE -package=mock -destination=mock/$GOFILE
package event

type Subscriber interface {
	Type() interface{}
	Subscribe(e Event) error
}
