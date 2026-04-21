package event

import (
	"encoding/asn1"
	"errors"

	"github.com/beatoz/bprn-sdk-go/chaincodes/event/merkle"
	"github.com/beatoz/bprn-sdk-go/chaincodes/event/types"
)

type eventLogHeader struct {
	ChannelId   string `json:"channel_id"`
	ChaincodeId string `json:"chaincode_id"`
	TxId        []byte `json:"tx_id"`
	Selector    []byte `json:"selector"`
}

func WithChannelId(channelId string) func(x *EventLog) {
	return func(x *EventLog) {
		x.Header.ChannelId = channelId
	}
}

func WithChaincodeId(ccId string) func(x *EventLog) {
	return func(x *EventLog) {
		x.Header.ChaincodeId = ccId
	}
}

func WithTxId(txId []byte) func(x *EventLog) {
	return func(x *EventLog) {
		x.Header.TxId = txId
	}
}

func WithSelector(selector []byte) func(x *EventLog) {
	return func(x *EventLog) {
		x.Header.Selector = selector
	}
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
		x.TxId,
		x.Selector,
	}
}

func (x *eventLogHeader) LeavesLen() int {
	return len(x.Leaves())
}

var _ types.ILeaves = (*eventLogHeader)(nil)

type EventLog struct {
	Header *eventLogHeader `json:"header"`
	Elems  [][]byte        `json:"elems"`
	tree   *merkle.MerkleTree
}

func NewEventLog(opt ...func(*EventLog)) *EventLog {
	evtlog := &EventLog{}
	for _, f := range opt {
		f(evtlog)
	}
	return evtlog
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

var _ types.ILeaves = (*EventLog)(nil)

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

var _ types.IMerkleProvable = (*EventLog)(nil)

func (log *EventLog) SetElems(elems types.IEventElems) {
	log.Header.Selector = elems.Selector()
	log.Elems = elems.Leaves()
	log.tree = nil
}

func (log *EventLog) Reset() {
	log.Elems = nil
	log.tree = nil
}

type derEventLog struct {
	ChannelId   string
	ChaincodeId string
	TxId        []byte
	Selector    []byte
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
		Selector:    log.Header.Selector,
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
		Selector:    d.Selector,
	}
	log.Elems = nil
	for _, raw := range d.Elems {
		log.Elems = append(log.Elems, raw.Bytes)
	}
	log.tree = nil
	return nil
}
