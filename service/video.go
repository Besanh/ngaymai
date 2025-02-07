package service

type (
	IVideo interface{}
	Video  struct{}
)

func NewVideo() IVideo {
	return &Video{}
}
