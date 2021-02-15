package channel

type Channel interface {
	Join(*Session, *Message) (interface{}, error)
	Handle(*Session, *Message) (interface{}, error)
}
