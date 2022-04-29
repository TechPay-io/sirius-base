package basestreamseeder

import (
	"bytes"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/TechPay-io/sirius-base/gossip/basestream"
	"github.com/TechPay-io/sirius-base/hash"
	"github.com/TechPay-io/sirius-base/inter/dag"
)

type testLocator struct {
	B []byte
}

func (l testLocator) Compare(b basestream.Locator) int {
	return bytes.Compare(l.B, b.(testLocator).B)
}

func (l testLocator) Inc() basestream.Locator {
	nextBn := new(big.Int).SetBytes(l.B)
	nextBn.Add(nextBn, common.Big1)
	return testLocator{
		B: nextBn.Bytes(),
	}
}

type testPayload struct {
	IDs    hash.Events
	Events dag.Events
	Size   uint64
}

func (p testPayload) AddEvent(id hash.Event, event dag.Event) {
	p.IDs = append(p.IDs, id)
	p.Events = append(p.Events, event)
	p.Size += uint64(event.Size())
}

func (p testPayload) AddID(id hash.Event) {
	p.IDs = append(p.IDs, id)
}

func (p testPayload) Len() int {
	return len(p.IDs)
}

func (p testPayload) TotalSize() uint64 {
	return p.Size
}

func (p testPayload) TotalMemSize() int {
	return int(p.Size) + len(p.IDs)*128
}
