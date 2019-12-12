package transition

import (
	"github.com/pkg/errors"
)

var (
	errTransitionExist       = errors.New("transition already exist")
	errTransitionInvalid     = errors.New("transition state not valid")
	errTransitionWrongStatus = errors.New("transition from current state not allowed")
)

const (
	NOOP string = "noop"
)

type Status int

type Transition struct {
	Name string
	Fn   func(ctx interface{}) error
	From Status
	To   Status
}

type Engine interface {
	Register(ts Transition) error
	Process(name string, ctx interface{}, current Status) (Status, error)
}

var _ Engine = &engine{}

type engine struct {
	states      []Status
	transitions map[string]Transition
}

func NewEngine(states []Status) Engine {
	return &engine{
		states:      states,
		transitions: make(map[string]Transition),
	}
}

func (m *engine) Register(ts Transition) error {
	if _, ok := m.transitions[ts.Name]; ok {
		return errTransitionExist
	}
	valid := 0
	for _, s := range m.states {
		if ts.From == s || ts.To == s {
			valid++
		}
	}
	if valid != 2 {
		return errTransitionInvalid
	}

	m.transitions[ts.Name] = ts
	return nil
}

func (m *engine) Process(name string, ctx interface{}, current Status) (Status, error) {

	if name == NOOP {
		return -1, nil
	}
	ts, ok := m.transitions[name]
	if !ok {
		panic("requested transition not registered: " + name)
	}
	if ts.From != current {
		return ts.From, errTransitionWrongStatus
	}
	err := ts.Fn(ctx)
	return ts.To, err
}
