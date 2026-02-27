package event

import (
	"github.com/beatoz/bprn-sdk-go/chaincodes/event/merkle"
	"github.com/holiman/uint256"
)

type TransferEventLog struct {
	From   []byte
	To     []byte
	Amount *uint256.Int
	Memo   []byte
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
