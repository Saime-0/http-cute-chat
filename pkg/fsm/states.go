package fsm

type State struct {
	Onset func() error
}

func NewState(onset func() error) *State {
	return &State{Onset: onset}
}

// with empty func
func NewTransitionalState() *State {
	return &State{Onset: func() error { return nil }}
}
