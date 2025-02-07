package repository

type (
	IChannel interface{}
	Channel  struct{}
)

func NewChannel() IChannel {
	return &Channel{}
}
