package merkle

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

// Test vector of BTIP-48 "Test Vectors / Tree".
func TestSpecVector_Tree(t *testing.T) {
	leaves := [][]byte{{0x61}, {0x62}, nil}
	tree := NewMerkleTree(WithRawLeaves(leaves))
	require.Equal(t, 4, tree.leafCount)

	want := map[int]string{
		1: "ed59f0e9e53a90b1fa76c23816424acaab94c5fd6073e597517daf43f6b38154",
		2: "b137985ff484fb600db93107c77b0365c80d78f5b429ded0fd97361d077999eb",
		3: "8c35feba66fbe78ac0ead640353127fa1bfe1c20c64b4ab6ce2ee56418828f0e",
		4: "022a6979e6dab7aa5ae4c3e5e45f7e977112a7e63593820dbec1ec738a24f93c",
		5: "57eb35615d47f34ec714cacdf5fd74608a5e8e102724e80b24b287c0c27b6a31",
		6: "6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d", // NULL_HASH
		7: "e3b0c44298fc1c149afbf4c8996fb92427" +
			"ae41e4649b934ca495991b7852b855", // EMPTY_HASH
	}
	for i, hexWant := range want {
		require.Equal(t, hexWant, hex.EncodeToString(tree.nodes[i]), "nodes[%d]", i)
	}
	require.Equal(t, want[1], hex.EncodeToString(tree.Root()))
}

// Test vector of BTIP-48 "Test Vectors / Proof".
func TestSpecVector_Proof(t *testing.T) {
	leaves := [][]byte{{0x61}, {0x62}, nil}
	tree := NewMerkleTree(WithRawLeaves(leaves))

	leaf, siblings, err := tree.Proof(1)
	require.NoError(t, err)
	require.Equal(t, []byte{0x62}, leaf)
	require.Len(t, siblings, 2)
	require.Equal(t, "022a6979e6dab7aa5ae4c3e5e45f7e977112a7e63593820dbec1ec738a24f93c",
		hex.EncodeToString(siblings[0]))
	require.Equal(t, "8c35feba66fbe78ac0ead640353127fa1bfe1c20c64b4ab6ce2ee56418828f0e",
		hex.EncodeToString(siblings[1]))

	require.NoError(t, VerifyProof(1, leaf, siblings, tree.Root()))
}

// TestNewMerkleTree_RawData tests tree construction with arbitrary leaf data.
func TestNewMerkleTree_RawData(t *testing.T) {
	leaves := [][]byte{
		[]byte("alice"),
		[]byte("bob"),
		[]byte("charlie"),
		[]byte("dave"),
	}
	tree := NewMerkleTree(WithRawLeaves(leaves))
	require.Len(t, tree.Root(), HashSize)

	want := InnerHash(
		InnerHash(LeafHash(leaves[0]), LeafHash(leaves[1])),
		InnerHash(LeafHash(leaves[2]), LeafHash(leaves[3])),
	)
	require.Equal(t, want, tree.Root())
}

// TestNewMerkleTree_PaddingToPowerOf2 tests that non-power-of-2 leaves are padded
// with EmptyHash, which is a domain of its own.
func TestNewMerkleTree_PaddingToPowerOf2(t *testing.T) {
	tree := NewMerkleTree(WithRawLeaves([][]byte{[]byte("a"), []byte("b"), []byte("c")}))
	require.Equal(t, 4, tree.leafCount)
	require.Equal(t, 3, tree.LeafCount(), "LeafCount reports original leaves")
	require.Equal(t, EmptyHash(), tree.nodes[tree.leafCount+3], "padding leaf should be EmptyHash")
	require.Len(t, tree.Root(), HashSize)
}

// TestNewMerkleTree_SingleLeaf tests the tree with a single leaf: no internal
// node is computed, so the root is the leaf hash itself.
func TestNewMerkleTree_SingleLeaf(t *testing.T) {
	leaf := []byte("only")
	tree := NewMerkleTree(WithRawLeaves([][]byte{leaf}))
	require.Equal(t, 1, tree.leafCount)
	require.Equal(t, LeafHash(leaf), tree.Root())
}

// TestNewMerkleTree_ZeroLeaves: an empty tree has a single empty leaf slot.
func TestNewMerkleTree_ZeroLeaves(t *testing.T) {
	require.Equal(t, EmptyHash(), NewMerkleTree(WithRawLeaves(nil)).Root())
	require.Equal(t, EmptyHash(), NewMerkleTree(WithRawLeaves([][]byte{})).Root())
}

// TestNewMerkleTree_NilLeafIsNotPadding: a slot holding nil is a leaf, a slot
// past the original leaves is not.
func TestNewMerkleTree_NilLeafIsNotPadding(t *testing.T) {
	tree := NewMerkleTree(WithRawLeaves([][]byte{[]byte("a"), nil}))
	require.Equal(t, NullHash(), tree.nodes[3])
	require.NotEqual(t, EmptyHash(), tree.nodes[3])

	// [a, nil] and [a] (padded with EmptyHash) must not share a root.
	require.NotEqual(t, NewMerkleTree(WithRawLeaves([][]byte{[]byte("a")})).Root(), tree.Root())
}

// TestProof_Valid tests that proof generation and verification work for every leaf.
func TestProof_Valid(t *testing.T) {
	leaves := [][]byte{
		[]byte("tx1"),
		[]byte("tx2"),
		[]byte("tx3"),
		[]byte("tx4"),
	}
	tree := NewMerkleTree(WithRawLeaves(leaves))

	for i, data := range leaves {
		leaf, siblings, err := tree.Proof(i)
		require.NoError(t, err)
		require.Equal(t, data, leaf)
		require.NoError(t, VerifyProof(i, leaf, siblings, tree.Root()))
	}
}

// TestProof_OutOfRange tests that Proof rejects invalid indices, padding included.
func TestProof_OutOfRange(t *testing.T) {
	tree := NewMerkleTree(WithRawLeaves([][]byte{[]byte("a"), []byte("b"), []byte("c")}))
	_, _, err := tree.Proof(-1)
	require.Error(t, err, "expected error for negative index")
	_, _, err = tree.Proof(3)
	require.Error(t, err, "expected error for the padding slot")
	_, _, err = tree.Proof(4)
	require.Error(t, err, "expected error for out of range index")
}

// TestVerify_WrongData tests that verification fails with incorrect data.
func TestVerify_WrongData(t *testing.T) {
	leaves := [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")}
	tree := NewMerkleTree(WithRawLeaves(leaves))
	_, siblings, err := tree.Proof(0)
	require.NoError(t, err)

	require.Error(t, VerifyProof(0, []byte("fake"), siblings, tree.Root()))
}

// TestVerify_WrongIndex tests that verification fails with an incorrect index.
func TestVerify_WrongIndex(t *testing.T) {
	leaves := [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")}
	tree := NewMerkleTree(WithRawLeaves(leaves))
	leaf, siblings, err := tree.Proof(0)
	require.NoError(t, err)

	require.Error(t, VerifyProof(1, leaf, siblings, tree.Root()))
}

// TestVerify_WrongRoot tests that verification fails against a different root.
func TestVerify_WrongRoot(t *testing.T) {
	leaves := [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")}
	tree := NewMerkleTree(WithRawLeaves(leaves))
	leaf, siblings, err := tree.Proof(0)
	require.NoError(t, err)

	require.Error(t, VerifyProof(0, leaf, siblings, LeafHash([]byte("fake root"))))
}

// TestVerify_TamperedProof tests that verification fails with a modified sibling.
func TestVerify_TamperedProof(t *testing.T) {
	leaves := [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")}
	tree := NewMerkleTree(WithRawLeaves(leaves))
	leaf, siblings, err := tree.Proof(0)
	require.NoError(t, err)

	siblings[0] = LeafHash([]byte("tampered"))
	require.Error(t, VerifyProof(0, leaf, siblings, tree.Root()))
}

// TestSecurity_InternalNodeAsLeaf: a tree whose leaves are the internal nodes of
// another tree must not reproduce that tree's root. Domain separation is what
// breaks the equality the pre-BTIP-48 implementation had.
func TestSecurity_InternalNodeAsLeaf(t *testing.T) {
	leaves := [][]byte{[]byte("A"), []byte("B"), []byte("C"), []byte("D")}
	tree := NewMerkleTree(WithRawLeaves(leaves))
	root := tree.Root()

	n1 := InnerHash(LeafHash(leaves[0]), LeafHash(leaves[1]))
	n2 := InnerHash(LeafHash(leaves[2]), LeafHash(leaves[3]))
	forgedTree := NewMerkleTree(WithRawLeaves([][]byte{n1, n2}))

	require.NotEqual(t, root, forgedTree.Root(), "VULNERABLE: internal nodes as leaves reproduce the root")

	forgedLeaf, forgedSiblings, err := forgedTree.Proof(0)
	require.NoError(t, err)
	require.Error(t, VerifyProof(0, forgedLeaf, forgedSiblings, root), "VULNERABLE: forged proof accepted")
}

// TestSecurity_ConcatenatedLeavesAsLeaf tests that concatenating two leaves
// into one cannot forge a valid proof.
func TestSecurity_ConcatenatedLeavesAsLeaf(t *testing.T) {
	leaves := [][]byte{[]byte("tx1"), []byte("tx2"), []byte("tx3"), []byte("tx4")}
	tree := NewMerkleTree(WithRawLeaves(leaves))

	concat := append([]byte("tx1"), []byte("tx2")...)
	forgedTree := NewMerkleTree(WithRawLeaves([][]byte{concat, []byte("tx3"), []byte("tx4")}))

	require.NotEqual(t, tree.Root(), forgedTree.Root(), "VULNERABLE: concatenated leaves produce the same root")
}

// TestSecurity_ForgedProofWithInternalNode: an internal node offered as leaf data
// is rejected, because the verifier always applies LeafHash to it.
func TestSecurity_ForgedProofWithInternalNode(t *testing.T) {
	leaves := [][]byte{[]byte("A"), []byte("B"), []byte("C"), []byte("D")}
	tree := NewMerkleTree(WithRawLeaves(leaves))
	root := tree.Root()

	n1 := InnerHash(LeafHash(leaves[0]), LeafHash(leaves[1]))
	_, siblings, err := tree.Proof(0)
	require.NoError(t, err)

	require.Error(t, VerifyProof(0, n1, siblings, root), "VULNERABLE: internal node accepted as leaf data")
}

// TestSecurity_ShorterTreeDifferentRoot: a shallower tree built from the internal
// nodes of a deeper one no longer shares its root, and its proofs are rejected.
func TestSecurity_ShorterTreeDifferentRoot(t *testing.T) {
	leaves := [][]byte{[]byte("A"), []byte("B"), []byte("C"), []byte("D")}
	tree4 := NewMerkleTree(WithRawLeaves(leaves))
	root := tree4.Root()

	n1 := InnerHash(LeafHash(leaves[0]), LeafHash(leaves[1]))
	n2 := InnerHash(LeafHash(leaves[2]), LeafHash(leaves[3]))
	tree2 := NewMerkleTree(WithRawLeaves([][]byte{n1, n2}))

	require.NotEqual(t, root, tree2.Root(), "VULNERABLE: shorter tree reproduces the root")

	for i, data := range leaves {
		idx := i % tree2.LeafCount()
		_, forgedSiblings, err := tree2.Proof(idx)
		require.NoError(t, err)
		require.Error(t, VerifyProof(idx, data, forgedSiblings, root),
			fmt.Sprintf("VULNERABLE: shorter tree proof accepted for leaf[%d]", i))
	}
}

// TestSubtreeComposition tests that a subtree root can be carried as a leaf of a
// parent tree. Under BTIP-48 it is ordinary leaf data there, so it is hashed
// again with LeafHash -- that is the binding between the two tree levels.
func TestSubtreeComposition(t *testing.T) {
	sub1 := NewMerkleTree(WithRawLeaves([][]byte{[]byte("a"), []byte("b")}))
	sub2 := NewMerkleTree(WithRawLeaves([][]byte{[]byte("c"), []byte("d")}))

	parent := NewMerkleTree(WithRawLeaves([][]byte{sub1.Root(), sub2.Root()}))
	require.Equal(t, InnerHash(LeafHash(sub1.Root()), LeafHash(sub2.Root())), parent.Root())

	leaf, siblings, err := parent.Proof(0)
	require.NoError(t, err)
	require.Equal(t, sub1.Root(), leaf)
	require.NoError(t, VerifyProof(0, leaf, siblings, parent.Root()))

	// The subtree root must not pass as a node of the parent tree itself.
	require.Error(t, VerifyProof(0, nil, nil, parent.Root()))
}

func TestNextPowerOf2(t *testing.T) {
	tests := []struct {
		input, expected int
	}{
		{0, 1}, {1, 1}, {2, 2}, {3, 4}, {4, 4}, {5, 8}, {7, 8}, {8, 8}, {9, 16},
	}
	for _, tc := range tests {
		got := nextPowerOf2(tc.input)
		require.Equal(t, tc.expected, got, fmt.Sprintf("nextPowerOf2(%d) = %d, want %d", tc.input, got, tc.expected))
	}
}
