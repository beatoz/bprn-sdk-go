package event

import (
	"fmt"
	"testing"

	"github.com/beatoz/bprn-sdk-go/chaincodes/event/merkle"
	"github.com/stretchr/testify/require"
)

var evtLog *EventLog

func init() {
	evtLog = NewEventLog("channelId", "chaincodeName", "txId")
	for i := 0; i < 10; i++ {
		logVal := &PostMessageEventLog{
			SrcChainId: fmt.Sprintf("srcChainId-%d", i),
			SrcDappId:  fmt.Sprintf("srcDappId-%d", i),
			SrcAcctId:  fmt.Sprintf("srcAcctId-%d", i),
			DstChainId: fmt.Sprintf("dstChainId-%d", i),
			DstDappId:  fmt.Sprintf("dstDappId-%d", i),
			DstAcctId:  fmt.Sprintf("dstAcctId-%d", i),
			MsgIdx:     uint64(i * 100),
			MsgPayload: []byte(fmt.Sprintf("hello world-%d", i)),
		}
		evtLog.AddLeaves(logVal)
	}
}

func TestCodec(t *testing.T) {
	bz, err := MarshalEventLog(evtLog)
	require.NoError(t, err)

	evtLog0, err := UnmarshalEventLog(bz)
	require.NoError(t, err)
	require.Equal(t, evtLog.Header.ChannelId, evtLog0.Header.ChannelId)
	require.Equal(t, evtLog.Header.ChaincodeName, evtLog0.Header.ChaincodeName)
	require.Equal(t, evtLog.Header.TxId, evtLog0.Header.TxId)
	require.Equal(t, len(evtLog.Elems), len(evtLog0.Elems))
	for i := 0; i < len(evtLog.Elems); i++ {
		for j := 0; j < len(evtLog.Elems[i]); j++ {
			require.Equal(t, evtLog.Elems[i][j], evtLog0.Elems[i][j])
		}
	}
}

func TestMerkleProof(t *testing.T) {
	evtLogRoot := evtLog.Root()

	gidx := 0 // channelId
	_, siblings, err := evtLog.Proof(gidx)
	require.NoError(t, err)
	require.NoError(t, evtLog.VerifyProof(gidx, siblings))
	require.NoError(t, merkle.VerifyProof(gidx, []byte("channelId"), siblings, evtLogRoot))

	gidx++
	_, siblings, err = evtLog.Proof(gidx)
	require.NoError(t, err)
	require.NoError(t, evtLog.VerifyProof(gidx, siblings))
	require.NoError(t, merkle.VerifyProof(gidx, []byte("chaincodeName"), siblings, evtLogRoot))

	gidx++
	_, siblings, err = evtLog.Proof(gidx)
	require.NoError(t, err)
	require.NoError(t, evtLog.VerifyProof(gidx, siblings))
	require.NoError(t, merkle.VerifyProof(gidx, []byte("txId"), siblings, evtLogRoot))

	leaf := "srcDappId-2" // header's leaves length(3) + PostMessageEventLog.LeavesLen(8)*2 + second field(1)
	gidx = 3 + 8*2 + 1
	_, siblings, err = evtLog.Proof(gidx)
	require.NoError(t, err)
	require.NoError(t, evtLog.VerifyProof(gidx, siblings))
	require.NoError(t, merkle.VerifyProof(gidx, []byte(leaf), siblings, evtLogRoot))
}
