package merkle

type ILeaves interface {
	LeavesLen() int
	Leaves() [][]byte
	Leaf(int) []byte
}
