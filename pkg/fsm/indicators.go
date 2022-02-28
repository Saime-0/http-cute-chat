package fsm

import (
	"github.com/saime-0/http-cute-chat/internal/cerrors"
	"time"
)

type (
	T      = interface{}
	Any    = interface{}
	States = map[T]*State
)

type Indicator struct {
	states    map[T]*State
	state     T
	ch        chan T
	listeners map[chan T]chan T
}

func (m *Machine) AddIndicator(name Any, defaultState T, states States) (*Indicator, error) {
	if len(states) < 2 {
		return nil, cerrors.New(StatesMapNotValid)
	}
	_, ok := states[defaultState]
	if !ok {
		return nil, cerrors.New(DefaultStateNotFound)
	}
	indicator := &Indicator{
		states: states,
		state:  defaultState,
		ch:     make(chan T),
	}
	m.Indicators[name] = indicator
	return indicator, nil
}

func (i *Indicator) SetState(key T) error {

	if key == i.state {
		return nil
	}

	state, ok := i.states[key]
	if !ok {
		return cerrors.New(StateNotFound)
	}

	err := state.Onset(i)
	if err != nil {
		return err
	}

	select {
	case i.ch <- state:
	default:
	}
	i.state = key

	return nil
}

func (i *Indicator) GetState() T {
	return i.state
}

func (i *Indicator) GetListener() <-chan T {
	ch := make(chan T)
	i.listeners[ch] = ch
	return ch
}

func (i *Indicator) loop(ch chan T) {
	for {
		<-ch
		time.Sleep(time.Millisecond)

		for listener := range i.listeners {
			select {
			case listener <- i.state:
			default:
				delete(i.listeners, listener)
				close(listener)
			}
		}
	}
}
