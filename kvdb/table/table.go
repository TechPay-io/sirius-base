package table

import (
	"github.com/Techpay-io/sirius-base/kvdb"
)

// Table wraper the underling DB, so all the table's data is stored with a prefix in underling DB
type Table struct {
	Readonly
	underlying kvdb.Store
}

var (
	// NOTE: key collisions are possible
	separator = []byte{}
)

// prefixed key (prefix + separator + key)
func prefixed(key, prefix []byte) []byte {
	prefixedKey := make([]byte, 0, len(prefix)+len(separator)+len(key))
	prefixedKey = append(prefixedKey, prefix...)
	prefixedKey = append(prefixedKey, separator...)
	prefixedKey = append(prefixedKey, key...)
	return prefixedKey
}

func noPrefix(key, prefix []byte) []byte {
	if len(key) < len(prefix)+len(separator) {
		return key
	}
	return key[len(prefix)+len(separator):]
}

/*
 * Database
 */

func New(db kvdb.Store, prefix []byte) *Table {
	return &Table{
		Readonly: Readonly{
			prefix:     prefix,
			underlying: db,
		},
		underlying: db,
	}
}

func (t *Table) NewTable(prefix []byte) *Table {
	return New(t, prefix)
}

func (t *Table) Close() error {
	return nil
}

// Drop the whole database.
func (t *Table) Drop() {}

func (t *Table) Put(key []byte, value []byte) error {
	return t.underlying.Put(prefixed(key, t.prefix), value)
}

func (t *Table) Delete(key []byte) error {
	return t.underlying.Delete(prefixed(key, t.prefix))
}

func (t *Table) NewBatch() kvdb.Batch {
	return &batch{t.underlying.NewBatch(), t.prefix}
}

func (t *Table) Compact(start []byte, limit []byte) error {
	return t.underlying.Compact(start, limit)
}

/*
 * Batch
 */

type batch struct {
	batch  kvdb.Batch
	prefix []byte
}

func (b *batch) Put(key, value []byte) error {
	return b.batch.Put(prefixed(key, b.prefix), value)
}

func (b *batch) Delete(key []byte) error {
	return b.batch.Delete(prefixed(key, b.prefix))
}

func (b *batch) ValueSize() int {
	return b.batch.ValueSize()
}

func (b *batch) Write() error {
	return b.batch.Write()
}

func (b *batch) Reset() {
	b.batch.Reset()
}

func (b *batch) Replay(w kvdb.Writer) error {
	return b.batch.Replay(&replayer{w, b.prefix})
}

/*
 * Replayer
 */

type replayer struct {
	writer kvdb.Writer
	prefix []byte
}

func (r *replayer) Put(key, value []byte) error {
	return r.writer.Put(noPrefix(key, r.prefix), value)
}

func (r *replayer) Delete(key []byte) error {
	return r.writer.Delete(noPrefix(key, r.prefix))
}
