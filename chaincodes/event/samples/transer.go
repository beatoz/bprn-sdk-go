package samples

import (
	"github.com/beatoz/bprn-sdk-go/chaincodes/event/merkle"
	"github.com/holiman/uint256"
)

type TransferEventLog struct {
	From   []byte       // gindex 3
	To     []byte       // gindex 4
	Amount *uint256.Int // gindex 5
	Memo   []byte       // gindex 6
}

func (s *TransferEventLog) Leaf(i int) []byte {
	if i >= s.LeavesLen() {
		return nil
	}
	return s.Leaves()[i]
}

func (s *TransferEventLog) Leaves() [][]byte {
	return [][]byte{
		s.From,
		s.To,
		s.Amount.Bytes(),
		s.Memo,
	}
}

func (s *TransferEventLog) LeavesLen() int {
	return len(s.Leaves())
}

var _ merkle.ILeaves = (*TransferEventLog)(nil)
