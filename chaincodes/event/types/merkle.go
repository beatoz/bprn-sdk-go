package types

type MerkleProof struct {
	Index    int      // 0-based position in the leaf array
	Leaf     []byte   // original leaf data
	Siblings [][]byte // sibling hashes, ordered from leaf level to root level
}

type IMerkleProvable interface {
	Root() []byte
	Proof(int) (*MerkleProof, error)
	VerifyProof(*MerkleProof) error
}

type ILeaves interface {
	Leaf(int) []byte
	Leaves() [][]byte
	LeavesLen() int
}
