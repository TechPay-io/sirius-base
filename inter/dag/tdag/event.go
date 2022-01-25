package tdag

import (
	"github.com/Techpay-io/sirius-base/hash"
	"github.com/Techpay-io/sirius-base/inter/dag"
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
