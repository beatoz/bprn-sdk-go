package merkle

import (
	"bytes"
	"errors"
	"fmt"
	"math/bits"

	"github.com/beatoz/bprn-sdk-go/chaincodes/event/types"
)

// MerkleTree is the array-based complete binary tree of BTIP-48.
//
// nodes is 1-indexed: nodes[1] is the root and the children of nodes[i] are
// nodes[2i] and nodes[2i+1]. Leaves occupy [leafCount, 2*leafCount-1], where
// leafCount is the original leaf count rounded up to a power of two; the slots
// past the original leaves hold EmptyHash.
type MerkleTree struct {
	nodes     [][]byte
	leaves    [][]byte // the original leaves as handed in, kept for proofs
	leafCount int      // n: leaf count after padding
}

type OptFunc func() [][]byte

func WithILeaves(leaves types.ILeaves) OptFunc {
	return func() [][]byte {
		return leaves.Leaves()
	}
}

func WithRawLeaves(leaves [][]byte) OptFunc {
	return func() [][]byte {
		return leaves
	}
}

func NewMerkleTree(opt OptFunc) *MerkleTree {
	return newMerkleTree(opt())
}

func newMerkleTree(leaves [][]byte) *MerkleTree {
	if len(leaves) > 1<<MaxMerkleDepth {
		return nil
	}

	leafCount := nextPowerOf2(len(leaves))
	tree := &MerkleTree{
		nodes:     make([][]byte, leafCount*2), // 1-indexed, nodes[0] is unused
		leaves:    leaves,
		leafCount: leafCount,
	}

	// populate leaves
	for i := 0; i < leafCount; i++ {
		switch {
		case i >= len(tree.leaves):
			tree.nodes[leafCount+i] = emptyHashAt[0] // padding slot: no leaf here
		case tree.leaves[i] == nil:
			tree.nodes[leafCount+i] = nullHashAt[0] // leaf slot holding nil
		default:
			tree.nodes[leafCount+i] = LeafHash(tree.leaves[i])
		}
	}

	height := 0 // height of the children of node i
	for i := leafCount - 1; i >= 1; i-- {
		left, right := tree.nodes[2*i], tree.nodes[2*i+1]
		switch {
		case bytes.Equal(left, emptyHashAt[height]) && bytes.Equal(right, emptyHashAt[height]):
			tree.nodes[i] = emptyHashAt[height+1]
		case bytes.Equal(left, nullHashAt[height]) && bytes.Equal(right, nullHashAt[height]):
			tree.nodes[i] = nullHashAt[height+1]
		default:
			tree.nodes[i] = InnerHash(left, right)
		}
		if i == leafCount>>(height+1) { // last node of this height
			height++
		}
	}

	return tree
}

// Root returns a copy of the merkle root hash.
func (t *MerkleTree) Root() []byte {
	if t == nil || len(t.nodes) < 2 {
		return nil
	}
	return cloneHash(t.nodes[1])
}

func (t *MerkleTree) LeafCount() int {
	if t == nil {
		return 0
	}
	return len(t.leaves)
}

// Proof returns the leaf data at index -- not its hash -- and its sibling
// hashes, leaf level first. The index is not returned; VerifyProof needs it.
func (t *MerkleTree) Proof(index int) ([]byte, [][]byte, error) {
	if t == nil {
		return nil, nil, errors.New("nil tree")
	}
	if index < 0 || index >= len(t.leaves) {
		return nil, nil, fmt.Errorf("index %d out of range [0, %d)", index, len(t.leaves))
	}

	siblings := make([][]byte, 0, bits.Len(uint(t.leafCount))-1)
	nodeIdx := t.leafCount + index
	for nodeIdx > 1 {
		// sibling is the XOR toggle of the last bit
		siblingIdx := nodeIdx ^ 1
		sibling := cloneHash(t.nodes[siblingIdx])
		siblings = append(siblings, sibling)
		nodeIdx /= 2 // move to parent
	}

	return cloneLeaf(t.leaves[index]), siblings, nil
}

// VerifyProof recomputes the root from the proof and compares it with root,
// which the caller must have obtained over a trusted path.
func VerifyProof(index int, leaf []byte, siblings [][]byte, root []byte) error {
	if err := ValidateProof(index, siblings); err != nil {
		return err
	}
	if len(root) != HashSize {
		return fmt.Errorf("invalid root length %d; expected %d", len(root), HashSize)
	}

	current := LeafHash(leaf)
	nodeIdx := index
	for _, sibling := range siblings {
		if nodeIdx%2 == 0 { // current is left child
			current = InnerHash(current, sibling)
		} else { // current is right child
			current = InnerHash(sibling, current)
		}
		nodeIdx /= 2
	}

	if !bytes.Equal(current, root) {
		return fmt.Errorf("root mismatch; computed %x, expected %x", current, root)
	}
	return nil
}

func ValidateProof(index int, siblings [][]byte) error {
	if len(siblings) > MaxMerkleDepth {
		return fmt.Errorf("proof too deep; max %d, got %d", MaxMerkleDepth, len(siblings))
	}

	for i, sibling := range siblings {
		if len(sibling) != HashSize {
			return fmt.Errorf("sibling %d: invalid length %d; expected %d", i, len(sibling), HashSize)
		}
	}

	n := 1 << len(siblings)
	if index < 0 || index >= n {
		return fmt.Errorf("index %d out of range [0, %d)", index, n)
	}
	return nil
}

func cloneLeaves(leaves [][]byte) [][]byte {
	if leaves == nil {
		return nil
	}
	out := make([][]byte, len(leaves))
	for i, leaf := range leaves {
		out[i] = cloneLeaf(leaf)
	}
	return out
}

func cloneLeaf(leaf []byte) []byte {
	if leaf == nil {
		return nil
	}
	return append([]byte{}, leaf...)
}

func cloneHash(hash []byte) []byte {
	return append([]byte(nil), hash...)
}

func nextPowerOf2(n int) int {
	if n <= 1 {
		return 1
	}
	return 1 << bits.Len(uint(n-1))
}
