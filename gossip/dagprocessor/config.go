package dagprocessor

import (
	"time"

	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/Techpay-io/sirius-base/inter/dag"
	"github.com/Techpay-io/sirius-base/utils/cachescale"
)

type Config struct {
	EventsBufferLimit dag.Metric

	EventsSemaphoreTimeout time.Duration

	MaxUnorderedInsertions int
}

func (c Config) MaxTasks() int {
	return c.MaxUnorderedInsertions*2 + 1
}

func DefaultConfig(scale cachescale.Func) Config {
	return Config{
		EventsBufferLimit: dag.Metric{
			// Shouldn't be too big because complexity is O(n) for each insertion in the EventsBuffer
			Num:  3000,
			Size: scale.U64(10 * opt.MiB),
		},
		EventsSemaphoreTimeout: 10 * time.Second,
	}
}
