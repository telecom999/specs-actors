package system

import (
	"github.com/filecoin-project/go-state-types/abi"

	"github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/runtime"
)

type Actor struct{}

func (a Actor) Exports() []interface{} {
	return []interface{}{
		builtin.MethodConstructor: a.Constructor,
	}
}

var _ runtime.Invokee = Actor{}

func (a Actor) Constructor(rt runtime.Runtime, _ *abi.EmptyValue) *abi.EmptyValue {
	rt.ValidateImmediateCallerIs(builtin.SystemActorAddr)

	rt.Create(&State{})
	return nil
}

type State struct{}
