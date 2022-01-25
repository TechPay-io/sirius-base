package peerleecher

import (
	"time"

	"github.com/Techpay-io/sirius-base/inter/dag"
)

type EpochDownloaderConfig struct {
	RecheckInterval        time.Duration
	DefaultChunkSize       dag.Metric
	ParallelChunksDownload int
}
