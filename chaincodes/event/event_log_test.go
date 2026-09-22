package event

import (
	"crypto/rand"
	"testing"

	"github.com/beatoz/bprn-sdk-go/chaincodes/event/merkle"
	"github.com/beatoz/bprn-sdk-go/chaincodes/event/samples"
	"github.com/stretchr/testify/require"
)

var evtLog *EventLog

var postMsgLog = &samples.PostMessageEventElems{
	SrcChainId: "srcChainId-0",
	SrcDappId:  "srcDappId-0",
	SrcAcctId:  "srcAcctId-0",
	DstChainId: "dstChainId-0",
	DstDappId:  "dstDappId-0",
	DstAcctId:  "dstAcctId-0",
	MsgIdx:     uint64(100),
	MsgPayload: []byte("hello world-0"),
}

var txId = make([]byte, 32)

func init() {
	_, _ = rand.Read(txId)
	evtLog = NewEventLog(
		WithChannelId("channelId"), WithChaincodeId("chaincodeName"), WithTxId(txId))
	evtLog.SetElems(postMsgLog)
}

func TestCodec(t *testing.T) {
	bz, err := evtLog.MarshalDER()
	require.NoError(t, err)

	evtLog0 := &EventLog{}
	err = evtLog0.UnmarshalDER(bz)
	require.NoError(t, err)
	require.Equal(t, evtLog.Header.ChannelId, evtLog0.Header.ChannelId)
	require.Equal(t, evtLog.Header.ChaincodeId, evtLog0.Header.ChaincodeId)
	require.Equal(t, evtLog.Header.TxId, evtLog0.Header.TxId)
	require.Equal(t, evtLog.Header.Selector, evtLog0.Header.Selector)
	require.Equal(t, evtLog.Elems, evtLog0.Elems)
	for i := 0; i < len(evtLog.Elems); i++ {
		for j := 0; j < len(evtLog.Elems[i]); j++ {
			require.Equal(t, evtLog.Elems[i][j], evtLog0.Elems[i][j])
		}
	}
	require.Equal(t, evtLog.Root(), evtLog0.Root())
}

func TestMerkleProof(t *testing.T) {
	evtLogRoot := evtLog.Root()

	for _, tc := range []struct {
		name string
		gidx int
		leaf []byte
	}{
		{"channelId", 0, []byte("channelId")},
		{"chaincodeId", 1, []byte("chaincodeName")},
		{"txId", 2, txId},
		{"selector", 3, postMsgLog.Selector()},
		{"srcDappId", 4 + 1, []byte("srcDappId-0")}, // header leaves(4) + second elem
	} {
		t.Run(tc.name, func(t *testing.T) {
			leaf, siblings, err := evtLog.Proof(tc.gidx)
			require.NoError(t, err)
			require.Equal(t, tc.leaf, leaf)

			require.NoError(t, evtLog.VerifyProof(tc.gidx, leaf, siblings))
			require.NoError(t, merkle.VerifyProof(tc.gidx, leaf, siblings, evtLogRoot))
		})
	}
}

func TestProof_IndexOutOfRange(t *testing.T) {
	_, _, err := evtLog.Proof(-1)
	require.Error(t, err)
	_, _, err = evtLog.Proof(evtLog.LeavesLen())
	require.Error(t, err)

	// The tree pads 12 leaves up to 16; the padding slots are not provable.
	require.Equal(t, 12, evtLog.LeavesLen())
	_, _, err = evtLog.Proof(12)
	require.Error(t, err)
}

func TestLeaf_IndexOutOfRange(t *testing.T) {
	require.Nil(t, evtLog.Leaf(-1))
	require.Nil(t, evtLog.Leaf(evtLog.LeavesLen()))
	require.Nil(t, evtLog.Header.Leaf(-1))
	require.Nil(t, evtLog.Header.Leaf(evtLog.Header.LeavesLen()))
}

// A proof of one event log must not verify against another log's root.
func TestVerifyProof_ForeignRoot(t *testing.T) {
	leaf, siblings, err := evtLog.Proof(0)
	require.NoError(t, err)

	other := NewEventLog(
		WithChannelId("otherChannel"), WithChaincodeId("chaincodeName"), WithTxId(txId))
	other.SetElems(postMsgLog)

	require.Error(t, merkle.VerifyProof(0, leaf, siblings, other.Root()))
	require.Error(t, other.VerifyProof(0, leaf, siblings))
}
