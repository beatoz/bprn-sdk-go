package samples

import (
	"encoding/binary"

	"github.com/beatoz/bprn-sdk-go/chaincodes/event/merkle"
)

type PostMessageEventLog struct {
	SrcChainId string `json:"src_chain_id,omitempty"`
	SrcDappId  string `json:"src_dapp_id,omitempty"`
	SrcAcctId  string `json:"src_acct_id,omitempty"`
	DstChainId string `json:"dst_chain_id,omitempty"`
	DstDappId  string `json:"dst_dapp_id,omitempty"`
	DstAcctId  string `json:"dst_acct_id,omitempty"`
	MsgIdx     uint64 `json:"msg_idx,omitempty"`
	MsgPayload []byte `json:"msg_payload,omitempty"`
}

func (p *PostMessageEventLog) Leaf(i int) []byte {
	if i >= p.LeavesLen() {
		return nil
	}
	return p.Leaves()[i]
}

func (p *PostMessageEventLog) Leaves() [][]byte {
	msgIdxBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(msgIdxBytes, p.MsgIdx)

	return [][]byte{
		[]byte(p.SrcChainId),
		[]byte(p.SrcDappId),
		[]byte(p.SrcAcctId),
		[]byte(p.DstChainId),
		[]byte(p.DstDappId),
		[]byte(p.DstAcctId),
		msgIdxBytes,
		p.MsgPayload,
	}
}

func (p *PostMessageEventLog) LeavesLen() int {
	return len(p.Leaves())
}

var _ merkle.ILeaves = (*PostMessageEventLog)(nil)
