package yubin

type Transport[Message any] interface {
	Send(message Message) error
}
