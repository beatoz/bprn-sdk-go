package event

import (
	"encoding/asn1"
	"errors"

	"github.com/beatoz/bprn-sdk-go/chaincodes/event/merkle"
)

type eventLogHeader struct {
	ChannelId   string `json:"channel_id"`
	ChaincodeId string `json:"chaincode_id"`
	TxId        string `json:"tx_id"`
}

func (x *eventLogHeader) Leaf(i int) []byte {
	if x.LeavesLen() <= i {
		return nil
	}
	return x.Leaves()[i]
}

func (x *eventLogHeader) Leaves() [][]byte {
	return [][]byte{
		[]byte(x.ChannelId),
		[]byte(x.ChaincodeId),
		[]byte(x.TxId),
	}
}

func (x *eventLogHeader) LeavesLen() int {
	return len(x.Leaves())
}

var _ merkle.ILeaves = (*eventLogHeader)(nil)

type EventLog struct {
	Header *eventLogHeader `json:"header"`
	Elems  [][]byte        `json:"elems"`
	tree   *merkle.MerkleTree
}

func NewEventLog(channelId, chaincodeId, txId string) *EventLog {
	return &EventLog{
		Header: &eventLogHeader{ChannelId: channelId, ChaincodeId: chaincodeId, TxId: txId},
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

type derEventLog struct {
	ChannelId   string
	ChaincodeId string
	TxId        string
	Elems       []asn1.RawValue
}

func (log *EventLog) MarshalDER(onlyElems ...bool) ([]byte, error) {
	if len(onlyElems) > 0 && onlyElems[0] {
		return asn1.Marshal(log.Elems)
	}
	d := derEventLog{
		ChannelId:   log.Header.ChannelId,
		ChaincodeId: log.Header.ChaincodeId,
		TxId:        log.Header.TxId,
	}
	for _, e := range log.Elems {
		d.Elems = append(d.Elems, asn1.RawValue{
			Class: asn1.ClassUniversal,
			Tag:   asn1.TagOctetString,
			Bytes: e,
		})
	}
	return asn1.Marshal(d)
}

func (log *EventLog) UnmarshalDER(data []byte, onlyElems ...bool) error {
	if len(onlyElems) > 0 && onlyElems[0] {
		var elems [][]byte
		_, err := asn1.Unmarshal(data, &elems)
		if err != nil {
			return err
		}
		log.Elems = elems
		log.tree = nil
		return nil
	}
	var d derEventLog
	_, err := asn1.Unmarshal(data, &d)
	if err != nil {
		return err
	}
	log.Header = &eventLogHeader{
		ChannelId:   d.ChannelId,
		ChaincodeId: d.ChaincodeId,
		TxId:        d.TxId,
	}
	log.Elems = nil
	for _, raw := range d.Elems {
		log.Elems = append(log.Elems, raw.Bytes)
	}
	log.tree = nil
	return nil
}
