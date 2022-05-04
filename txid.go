package matrix

import (
	"fmt"
	"sync/atomic"
	"time"
)

// TXIDGen is responsible for generating a unique transaction ID that is used when creating events in matrix.
// this ID needs to be unique across all senders of this device. To conform to that you must only have one
// ID generator per device!
type TXIDGen interface {
	// NextTXID returns the next unique transaction ID as string.
	NextTXID() string
}

type txIDGen struct{ txID int64 }

func (t *txIDGen) NextTXID() string { return fmt.Sprint(atomic.AddInt64(&t.txID, 1)) }

// NewTXIDGen returns a new transaction ID generator.
func NewTXIDGen() TXIDGen { return &txIDGen{time.Now().UnixMilli()} }
