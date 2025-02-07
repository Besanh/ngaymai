package repository

type (
	IVideoRanking interface{}
	VideoRanking  struct{}
)

func NewVideoRanking() IVideoRanking {
	return &VideoRanking{}
}
