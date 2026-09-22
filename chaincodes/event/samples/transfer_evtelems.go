package samples

import (
	"crypto/sha256"

	"github.com/beatoz/bprn-sdk-go/chaincodes/event/types"
	"github.com/holiman/uint256"
)

// TransferEventElems implements BTIP25.
type TransferEventElems struct {
	From   []byte       // gidx 4
	To     []byte       // gidx 5
	Amount *uint256.Int // gidx 6
	Memo   []byte       // gidx 7
}

func (s *TransferEventElems) Leaf(i int) []byte {
	if i >= s.LeavesLen() {
		return nil
	}
	return s.Leaves()[i]
}

func (s *TransferEventElems) Leaves() [][]byte {
	return [][]byte{
		s.From,
		s.To,
		s.Amount.Bytes(),
		s.Memo,
	}
}

func (s *TransferEventElems) LeavesLen() int {
	return len(s.Leaves())
}

func (s *TransferEventElems) Selector() []byte {
	h := sha256.Sum256([]byte("TransferEventElems(bytes,bytes,uint256,bytes)"))
	return h[:]
}

var _ types.IEventElems = (*TransferEventElems)(nil)
