package fsm

type State struct {
	Onset func(indicator *Indicator) error
}

func NewState(onset func(indicator *Indicator) error) *State {
	return &State{Onset: onset}
}

// with empty func
func NewTransitionalState() *State {
	return &State{Onset: func(indicator *Indicator) error { return nil }}
}
