package abft

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/TechPay-io/sirius-base/abft/dagidx"
	"github.com/TechPay-io/sirius-base/hash"
	"github.com/TechPay-io/sirius-base/inter/dag"
	"github.com/TechPay-io/sirius-base/inter/idx"
	"github.com/TechPay-io/sirius-base/inter/pos"
	"github.com/TechPay-io/sirius-base/kvdb"
	"github.com/TechPay-io/sirius-base/sirius"
)

var _ sirius.Consensus = (*IndexedSirius)(nil)

// IndexedSirius performs events ordering and detects cheaters
// It's a wrapper around Orderer, which adds features which might potentially be application-specific:
// confirmed events traversal, DAG index updates and cheaters detection.
// Use this structure if need a general-purpose consensus. Instead, use lower-level abft.Orderer.
type IndexedSirius struct {
	*Sirius
	dagIndexer    DagIndexer
	uniqueDirtyID uniqueID
}

type DagIndexer interface {
	dagidx.VectorClock
	dagidx.ForklessCause

	Add(dag.Event) error
	Flush()
	DropNotFlushed()

	Reset(validators *pos.Validators, db kvdb.Store, getEvent func(hash.Event) dag.Event)
}

// New creates IndexedSirius instance.
func NewIndexedSirius(store *Store, input EventSource, dagIndexer DagIndexer, crit func(error), config Config) *IndexedSirius {
	p := &IndexedSirius{
		Sirius:        NewSirius(store, input, dagIndexer, crit, config),
		dagIndexer:    dagIndexer,
		uniqueDirtyID: uniqueID{new(big.Int)},
	}

	return p
}

// Build fills consensus-related fields: Frame, IsRoot
// returns error if event should be dropped
func (p *IndexedSirius) Build(e dag.MutableEvent) error {
	e.SetID(p.uniqueDirtyID.sample())

	defer p.dagIndexer.DropNotFlushed()
	err := p.dagIndexer.Add(e)
	if err != nil {
		return err
	}

	return p.Sirius.Build(e)
}

// Process takes event into processing.
// Event order matter: parents first.
// All the event checkers must be launched.
// Process is not safe for concurrent use.
func (p *IndexedSirius) Process(e dag.Event) (err error) {
	defer p.dagIndexer.DropNotFlushed()
	err = p.dagIndexer.Add(e)
	if err != nil {
		return err
	}

	err = p.Sirius.Process(e)
	if err != nil {
		return err
	}
	p.dagIndexer.Flush()
	return nil
}

func (p *IndexedSirius) Bootstrap(callback sirius.ConsensusCallbacks) error {
	base := p.Sirius.OrdererCallbacks()
	ordererCallbacks := OrdererCallbacks{
		ApplyAtropos: base.ApplyAtropos,
		EpochDBLoaded: func(epoch idx.Epoch) {
			if base.EpochDBLoaded != nil {
				base.EpochDBLoaded(epoch)
			}
			p.dagIndexer.Reset(p.store.GetValidators(), p.store.epochTable.VectorIndex, p.input.GetEvent)
		},
	}
	return p.Sirius.BootstrapWithOrderer(callback, ordererCallbacks)
}

type uniqueID struct {
	counter *big.Int
}

func (u *uniqueID) sample() [24]byte {
	u.counter = u.counter.Add(u.counter, common.Big1)
	var id [24]byte
	copy(id[:], u.counter.Bytes())
	return id
}
