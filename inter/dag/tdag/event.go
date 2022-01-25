package tdag

import (
	"github.com/Techpay-foundation/sirius-base/hash"
	"github.com/Techpay-foundation/sirius-base/inter/dag"
)

type TestEvent struct {
	dag.MutableBaseEvent
	Name string
}

func (e *TestEvent) AddParent(id hash.Event) {
	parents := e.Parents()
	parents.Add(id)
	e.SetParents(parents)
}
