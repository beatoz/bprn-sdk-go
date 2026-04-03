package event

import (
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

func init() {
	evtLog = NewEventLog("channelId", "chaincodeName", "txId")
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

	gidx := 0 // channelId
	_, siblings, err := evtLog.Proof(gidx)
	require.NoError(t, err)
	require.NoError(t, evtLog.VerifyProof(gidx, siblings))
	require.NoError(t, merkle.VerifyProof(gidx, []byte("channelId"), siblings, evtLogRoot))

	gidx++ // chaincodeId
	_, siblings, err = evtLog.Proof(gidx)
	require.NoError(t, err)
	require.NoError(t, evtLog.VerifyProof(gidx, siblings))
	require.NoError(t, merkle.VerifyProof(gidx, []byte("chaincodeName"), siblings, evtLogRoot))

	gidx++ // txId
	_, siblings, err = evtLog.Proof(gidx)
	require.NoError(t, err)
	require.NoError(t, evtLog.VerifyProof(gidx, siblings))
	require.NoError(t, merkle.VerifyProof(gidx, []byte("txId"), siblings, evtLogRoot))

	gidx++ // selector
	_, siblings, err = evtLog.Proof(gidx)
	require.NoError(t, err)
	require.NoError(t, evtLog.VerifyProof(gidx, siblings))
	require.NoError(t, merkle.VerifyProof(gidx, postMsgLog.Selector(), siblings, evtLogRoot))

	leaf := "srcDappId-0" // header's leaves length(4) + second field(1)
	gidx = 4 + 1
	_, siblings, err = evtLog.Proof(gidx)
	require.NoError(t, err)
	require.NoError(t, evtLog.VerifyProof(gidx, siblings))
	require.NoError(t, merkle.VerifyProof(gidx, []byte(leaf), siblings, evtLogRoot))
}
