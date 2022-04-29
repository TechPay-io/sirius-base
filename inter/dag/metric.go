package dag

import (
	"fmt"

	"github.com/TechPay-io/sirius-base/inter/idx"
)

type Metric struct {
	Num  idx.Event
	Size uint64
}

func (m Metric) String() string {
	return fmt.Sprintf("{Num=%d,Size=%d}", m.Num, m.Size)
}
