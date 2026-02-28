package event

import (
	"errors"

	"github.com/beatoz/bprn-sdk-go/chaincodes/event/merkle"
	"google.golang.org/protobuf/proto"
)

// EventLogHeaderProto
func (x *EventLogHeaderProto) Leaf(i int) []byte {
	if x.LeavesLen() <= i {
		return nil
	}
	return x.Leaves()[i]
}

func (x *EventLogHeaderProto) Leaves() [][]byte {
	return [][]byte{
		[]byte(x.ChannelId),
		[]byte(x.ChaincodeId),
		[]byte(x.TxId),
	}
}

func (x *EventLogHeaderProto) LeavesLen() int {
	return len(x.Leaves())
}

type EventLog struct {
	*EventLogProto
	tree *merkle.MerkleTree
}

func NewEventLog(channelId, chaincodeId, txId string) *EventLog {
	return &EventLog{
		EventLogProto: &EventLogProto{Header: &EventLogHeaderProto{ChannelId: channelId, ChaincodeId: chaincodeId, TxId: txId}},
	}
}

func (log *EventLog) Leaf(gidx int) []byte {
	if gidx >= log.LeavesLen() {
		return nil
	}
	if gidx < log.Header.LeavesLen() {
		return log.Header.Leaf(gidx)
	}
	gidx -= log.Header.LeavesLen()
	return log.Elems[gidx]
}

func (log *EventLog) Leaves() [][]byte {
	return append(log.Header.Leaves(), log.Elems...)
}

func (log *EventLog) LeavesLen() int {
	return log.Header.LeavesLen() + len(log.Elems)
}

var _ merkle.ILeaves = (*EventLog)(nil)

func (log *EventLog) Root() []byte {
	if log.tree == nil {
		log.tree = merkle.NewMerkleTree(merkle.WithILeaves(log))
	}
	return log.tree.Root()
}

func (log *EventLog) Proof(gidx int) ([]byte, [][]byte, error) {
	if gidx >= log.LeavesLen() {
		return nil, nil, errors.New("index out of range")
	}
	if log.tree == nil {
		log.tree = merkle.NewMerkleTree(merkle.WithILeaves(log))
	}
	return log.tree.Proof(gidx)
}

func (log *EventLog) VerifyProof(gidx int, siblings [][]byte) error {
	if gidx >= log.LeavesLen() {
		return errors.New("index out of range")
	}
	if log.tree == nil {
		log.tree = merkle.NewMerkleTree(merkle.WithILeaves(log))
	}
	return merkle.VerifyProof(gidx, log.Leaves()[gidx], siblings, log.tree.Root())
}

var _ merkle.IMerkleProvable = (*EventLog)(nil)

func (log *EventLog) AddData(data []byte) {
	log.Elems = append(log.Elems, data)
	log.tree = nil
}

func (log *EventLog) AddLeaves(logs merkle.ILeaves) {
	log.Elems = append(log.Elems, logs.Leaves()...)
	log.tree = nil
}

func (log *EventLog) SetLogs(logs merkle.ILeaves) {
	log.Elems = logs.Leaves()
	log.tree = nil
}

func (log *EventLog) Reset() {
	log.Elems = nil
	log.tree = nil
}

func MarshalEventLog(log *EventLog) ([]byte, error) {
	return proto.Marshal(log.EventLogProto)
}

func UnmarshalEventLog(data []byte) (*EventLog, error) {
	unmarshalled := &EventLog{
		EventLogProto: &EventLogProto{},
	}
	if err := proto.Unmarshal(data, unmarshalled.EventLogProto); err != nil {
		return nil, err
	}
	return unmarshalled, nil
}
