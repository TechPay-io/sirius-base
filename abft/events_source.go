package abft

import (
	"github.com/TechPay-io/sirius-base/hash"
	"github.com/TechPay-io/sirius-base/inter/dag"
)

// EventSource is a callback for getting events from an external storage.
type EventSource interface {
	HasEvent(hash.Event) bool
	GetEvent(hash.Event) dag.Event
}
