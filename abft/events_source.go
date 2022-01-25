package abft

import (
	"github.com/Techpay-foundation/sirius-base/hash"
	"github.com/Techpay-foundation/sirius-base/inter/dag"
)

// EventSource is a callback for getting events from an external storage.
type EventSource interface {
	HasEvent(hash.Event) bool
	GetEvent(hash.Event) dag.Event
}
