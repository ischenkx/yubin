package yubin

type Delivery[Message any] interface {
	Deliver(message Message) error
}
