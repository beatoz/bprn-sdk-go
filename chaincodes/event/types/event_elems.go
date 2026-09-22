package types

type IEventElems interface {
	ILeaves
	Selector() []byte
}
