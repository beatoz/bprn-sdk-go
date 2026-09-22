package merkle

import "crypto/sha256"

// Node hashing rules of BTIP-48.

// empty : zero-length input
// leaf  : 0x00 || data        (1 byte or more)
// inner : 0x01 || left||right (exactly 65 bytes)
const (
	LeafPrefix  byte = 0x00
	InnerPrefix byte = 0x01
)

const HashSize = sha256.Size

// MaxMerkleDepth caps proof depth so 1<<depth stays a positive int, even at 32 bits.
const MaxMerkleDepth = 30

// Per-height hashes of all-empty (padding) and all-null subtrees. Never mutated.
var (
	emptyHashAt [][]byte
	nullHashAt  [][]byte
)

func init() {
	emptyHashAt = make([][]byte, MaxMerkleDepth+1)
	nullHashAt = make([][]byte, MaxMerkleDepth+1)

	empty := sha256.Sum256(nil)
	emptyHashAt[0] = empty[:]
	nullHashAt[0] = LeafHash(nil)

	for h := 1; h <= MaxMerkleDepth; h++ {
		emptyHashAt[h] = InnerHash(emptyHashAt[h-1], emptyHashAt[h-1])
		nullHashAt[h] = InnerHash(nullHashAt[h-1], nullHashAt[h-1])
	}
}

// LeafHash returns sha256(0x00 || data). nil and an empty slice hash alike.
func LeafHash(data []byte) []byte {
	h := sha256.New()
	h.Write([]byte{LeafPrefix})
	h.Write(data)
	return h.Sum(nil)
}

// InnerHash returns sha256(0x01 || left || right), or nil unless both children are 32 bytes.
func InnerHash(left, right []byte) []byte {
	if len(left) != HashSize || len(right) != HashSize {
		return nil
	}
	h := sha256.New()
	h.Write([]byte{InnerPrefix})
	h.Write(left)
	h.Write(right)
	return h.Sum(nil)
}

func EmptyHash() []byte { return EmptyHashAt(0) }

func NullHash() []byte { return NullHashAt(0) }

func EmptyHashAt(height int) []byte {
	if height < 0 || height > MaxMerkleDepth {
		return nil
	}
	return append([]byte(nil), emptyHashAt[height]...)
}

func NullHashAt(height int) []byte {
	if height < 0 || height > MaxMerkleDepth {
		return nil
	}
	return append([]byte(nil), nullHashAt[height]...)
}
