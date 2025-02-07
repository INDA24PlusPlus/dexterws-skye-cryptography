package utils

type u512 [8]uint64
type MerkleNode struct {
	Data  u512
	Left  *MerkleNode
	Right *MerkleNode
}
type MerkleTree struct {
	Root      *MerkleNode
	ValidHash u512
}

func u512eq(a u512, b u512) (eq bool) {
	eq = true
	for i := 0; i < 8; i++ {
		eq = eq && (a[i] == b[i])
	}
	return
}
func u512bwxor(a u512, b u512) (c u512) {
	for i := 0; i < 8; i++ {
		c[i] = a[i] ^ b[i]
	}
	return
}
func Hash(h1 u512, h2 u512) (hash u512) {
	hash = u512bwxor(h1, h2)
	return
}
func (MN *MerkleNode) ValidateLeaf(path uint8, depth uint8, check u512) (hash u512) {
	if depth <= 7 {
		if ((0x01 << (7 - depth)) & path) != 0x00 {
			hash = Hash(MN.Right.ValidateLeaf(path, depth+1, check), MN.Left.Data)
		} else {
			hash = Hash(MN.Left.ValidateLeaf(path, depth+1, check), MN.Right.Data)
		}
		return
	}
	hash = check
	return
}

func (MN *MerkleNode) CalcHash() (hash u512) {
	if MN.Left != nil {
		MN.Data = Hash(MN.Left.CalcHash(), MN.Right.CalcHash())
		hash = MN.Data
		return
	}
	hash = MN.Data
	return
}

func (MT *MerkleTree) CalcHash() (hash u512) {
	MT.ValidHash = MT.Root.CalcHash()
	hash = MT.ValidHash
	return
}
func (MT *MerkleTree) ValidateHash(path uint8, check u512) (valid bool) {
	valid = u512eq(MT.Root.ValidateLeaf(path, 0, check), MT.ValidHash)
	return
}
