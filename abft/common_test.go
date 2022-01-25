package abft

import (
	"github.com/TechPay-io/sirius-base/inter/idx"
	"github.com/TechPay-io/sirius-base/inter/pos"
	"github.com/TechPay-io/sirius-base/kvdb"
	"github.com/TechPay-io/sirius-base/kvdb/memorydb"
	"github.com/TechPay-io/sirius-base/sirius"
	"github.com/TechPay-io/sirius-base/utils/adapters"
	"github.com/TechPay-io/sirius-base/vecfc"
)

type applyBlockFn func(block *sirius.Block) *pos.Validators

// TestSirius extends Sirius for tests.
type TestSirius struct {
	*IndexedSirius

	blocks map[idx.Block]*sirius.Block

	applyBlock applyBlockFn
}

// FakeSirius creates empty abft with mem store and equal weights of nodes in genesis.
func FakeSirius(nodes []idx.ValidatorID, weights []pos.Weight, mods ...memorydb.Mod) (*TestSirius, *Store, *EventStore) {
	validators := make(pos.ValidatorsBuilder, len(nodes))
	for i, v := range nodes {
		if weights == nil {
			validators[v] = 1
		} else {
			validators[v] = weights[i]
		}
	}

	openEDB := func(epoch idx.Epoch) kvdb.DropableStore {
		return memorydb.New()
	}
	crit := func(err error) {
		panic(err)
	}
	store := NewStore(memorydb.New(), openEDB, crit, LiteStoreConfig())

	err := store.ApplyGenesis(&Genesis{
		Validators: validators.Build(),
		Epoch:      FirstEpoch,
	})
	if err != nil {
		panic(err)
	}

	input := NewEventStore()

	config := LiteConfig()
	lch := NewIndexedSirius(store, input, &adapters.VectorToDagIndexer{vecfc.NewIndex(crit, vecfc.LiteConfig())}, crit, config)

	extended := &TestSirius{
		IndexedSirius: lch,
		blocks:        map[idx.Block]*sirius.Block{},
	}

	blockIdx := idx.Block(0)

	err = extended.Bootstrap(sirius.ConsensusCallbacks{
		BeginBlock: func(block *sirius.Block) sirius.BlockCallbacks {
			blockIdx++
			return sirius.BlockCallbacks{
				EndBlock: func() (sealEpoch *pos.Validators) {
					// track blocks
					extended.blocks[blockIdx] = block
					if extended.applyBlock != nil {
						return extended.applyBlock(block)
					}
					return nil
				},
			}
		},
	})
	if err != nil {
		panic(err)
	}

	return extended, store, input
}
