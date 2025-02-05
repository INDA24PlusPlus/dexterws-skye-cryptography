package main

type u512 [8]uint64
type MerkleNode struct {
	data  u512
	Left  *MerkleNode
	Right *MerkleNode
}
type MerkleTree struct {
	root *MerkleNode
}
func u512bwxor(a u512, b u512) (c u512){
	for i:=0; i<8; i++ {
		c[i]=a[i]^b[i]
	}
	return
}
func Hash(h1 u512, h2 u512) (hash u512){
	hash=u512bwxor(h1,h2)
	return
}

func (MN *MerkleNode) CalcHash() (hash u512){
	if(MN.Left!=nil){
		MN.data=Hash(MN.Left.CalcHash(),MN.Right.CalcHash())
		hash=MN.data
		return
	}
	hash=MN.data
	return
}
func (MT *MerkleTree) GetHash() (hash u512) {
	hash=MT.root.CalcHash()
	return
}

func main(){}
