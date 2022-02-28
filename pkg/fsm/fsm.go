package fsm

type Machine struct {
	Indicators map[interface{}]*Indicator
}

func NewMachine() (*Machine, error) {
	//if len(indicators) == 0 {
	//	return nil, cerrors.New(NoOneIndicatorNotCreated)
	//}
	return &Machine{
		Indicators: map[interface{}]*Indicator{},
	}, nil
}
