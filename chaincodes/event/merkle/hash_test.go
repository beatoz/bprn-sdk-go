package merkle

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test vectors of BTIP-48 "Test Vectors / Constants".
func TestSpecVector_Constants(t *testing.T) {
	emptyWant := []string{
		"e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
		"a68ee79dc12813d134fd035c7328f7bd5ee68187735f7f0d2e451aea3ff6930f",
		"905e5783421ee97400a05d92dae5289a5d689e527e8d1b3fba314226e725f3e5",
		"57fda7f49b727e354d82f0c79adc240e0e3edaa70d53426cd890fd9cf713e657",
	}
	nullWant := []string{
		"6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d",
		"fe43d66afa4a9a5c4f9c9da89f4ffb52635c8f342e7ffb731d68e36c5982072a",
		"deb82e155954d6be14592c66ccf7a1ece193eeebcdabaf747b91f44519f09f47",
		"2960044c62f2354e945e8d78fdd220a05f2c0879f24df6f11ef5cc26b5270a0e",
	}

	for h, want := range emptyWant {
		require.Equal(t, want, hex.EncodeToString(EmptyHashAt(h)), "EMPTY_HASH_AT[%d]", h)
	}
	for h, want := range nullWant {
		require.Equal(t, want, hex.EncodeToString(NullHashAt(h)), "NULL_HASH_AT[%d]", h)
	}

	require.Equal(t, EmptyHashAt(0), EmptyHash())
	require.Equal(t, NullHashAt(0), NullHash())
}

// EMPTY_HASH is sha256 over an empty input, NULL_HASH is a leaf holding nil.
// The two must not collide, otherwise a padding slot could be proven as a leaf.
func TestDomainSeparation_NullVsEmpty(t *testing.T) {
	sum := sha256.Sum256(nil)
	require.Equal(t, sum[:], EmptyHash(), "EmptyHash should be sha256(nil)")
	require.Equal(t, LeafHash(nil), NullHash(), "NullHash should be LeafHash(nil)")
	require.NotEqual(t, EmptyHash(), NullHash())
}

func TestLeafHash_NilEqualsEmptySlice(t *testing.T) {
	require.Equal(t, LeafHash(nil), LeafHash([]byte{}))
}

// A leaf hash and an internal node hash must never be produced from the same
// input, which is what lets a verifier reject an internal node offered as a leaf.
func TestDomainSeparation_LeafVsInner(t *testing.T) {
	left, right := LeafHash([]byte("a")), LeafHash([]byte("b"))
	inner := InnerHash(left, right)

	require.Len(t, inner, HashSize)
	require.NotEqual(t, inner, LeafHash(append(append([]byte{}, left...), right...)),
		"inner hash must differ from a leaf holding the same 64 bytes")
}

func TestInnerHash_RejectsMalformedChildren(t *testing.T) {
	ok := LeafHash([]byte("x"))
	require.Nil(t, InnerHash(ok[:31], ok))
	require.Nil(t, InnerHash(ok, ok[:31]))
	require.Nil(t, InnerHash(nil, ok))
	require.Nil(t, InnerHash(ok, append(append([]byte{}, ok...), 0x00)))
	require.NotNil(t, InnerHash(ok, ok))
}

func TestHashAt_OutOfRange(t *testing.T) {
	require.Nil(t, EmptyHashAt(-1))
	require.Nil(t, NullHashAt(-1))
	require.Nil(t, EmptyHashAt(MaxMerkleDepth+1))
	require.Nil(t, NullHashAt(MaxMerkleDepth+1))
	require.NotNil(t, EmptyHashAt(MaxMerkleDepth))
	require.NotNil(t, NullHashAt(MaxMerkleDepth))
}

// The tables are package state; handing out the backing array would let one
// caller corrupt every later tree.
func TestHashAt_ReturnsCopy(t *testing.T) {
	got := EmptyHashAt(1)
	got[0] ^= 0xff
	require.NotEqual(t, got, EmptyHashAt(1), "EmptyHashAt must return a copy")

	got = NullHashAt(1)
	got[0] ^= 0xff
	require.NotEqual(t, got, NullHashAt(1), "NullHashAt must return a copy")
}

// The height tables must equal the plain recursive definition.
func TestHashAt_MatchesRecursiveDefinition(t *testing.T) {
	empty, null := EmptyHash(), NullHash()
	for h := 1; h <= MaxMerkleDepth; h++ {
		empty = InnerHash(empty, empty)
		null = InnerHash(null, null)
		require.Equal(t, empty, EmptyHashAt(h), "EMPTY_HASH_AT[%d]", h)
		require.Equal(t, null, NullHashAt(h), "NULL_HASH_AT[%d]", h)
	}
}
