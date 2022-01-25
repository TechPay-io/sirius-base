package fallible

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Fantom-foundation/lachesis-base/kvdb"
	"github.com/Fantom-foundation/lachesis-base/kvdb/memorydb"
)

func TestFallible(t *testing.T) {
	assertar := assert.New(t)

	var (
		key = []byte("test-key")
		val = []byte("test-value")
		db  kvdb.Store
		err error
	)

	mem := memorydb.New()
	w := Wrap(mem)
	db = w

	_, err = db.Get(key)
	assertar.NoError(err)

	assertar.Panics(func() {
		db.Put(key, val)
	})

	w.SetWriteCount(1)

	err = db.Put(key, val)
	assertar.NoError(err)

	assertar.Panics(func() {
		err = db.Put(key, val)
	})
}
