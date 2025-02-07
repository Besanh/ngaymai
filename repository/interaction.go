package repository

type (
	IInteraction interface{}
	Interaction  struct{}
)

func NewInteraction() IInteraction {
	return &Interaction{}
}
