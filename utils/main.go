package utils
/*
import (
	"encoding/binary"
	"crypto/sha512"
	"fmt"
)
type u512 [8]uint64
type b512 [64]byte;
func add_b512(a b512, b b512) (c b512) {
	for i:=0; i<64; i++{
		c[i]=a[i]+b[i]
	}
	return
}
func b512_to_bytes(b b512) (o []byte){
	for i:=0; i<64; i++{
		o[i]=b[i]
	}
	return
}
func bytes_to_b512(b [64]byte) (o b512){
	for i:=0; i<64; i++{
		o[i]=b[i]
	}
	return
}

func to_b512(u u512) (b b512){
	 for i:=0; i<8; i++{
	 	tB := make([]byte, 8)
	 	binary.LittleEndian.PutUint64(tB, u[i])
		for j:=0; j<8; j++ {
			b[i*8+j]=tB[j]
		}

	 }
	 return
}
type MerkleNode struct {
	Data  	b512
	Left  	*MerkleNode
	Right 	*MerkleNode
}
type MerkleTree struct {
	Root *MerkleNode
	ValidHash b512
}
func b512eq(a b512, b b512) (eq bool){
	eq=true
	for i:=0; i<64; i++ {
		eq=eq&&(a[i]==b[i])
	}
	return
}
func b512bwxor(a b512, b b512) (c b512){
	for i:=0; i<64; i++ {
		c[i]=a[i]^b[i]
	}
	return
}
func Hash(h1 b512, h2 b512) (hash b512){
	hash=sha512.Sum512(b512_to_bytes(add_b512(h1,h2)))
	return
}

func (MN *MerkleNode) ValidateLeaf(path uint8, depth uint8, check b512) (hash b512){
	if(depth<=7){
		if(((0x01<<(depth))&path)!=0x00) {
			hash=Hash(MN.Right.ValidateLeaf(path, depth+1, check),MN.Left.Data)
		} else {
			hash=Hash(MN.Left.ValidateLeaf(path, depth+1, check), MN.Right.Data)
		}
		return
	}
	hash=check
	return
}
func (MN *MerkleNode) SetLeaf(path uint8, depth uint8, check b512){
	if(depth<=7){
		if(((0x01<<(depth))&path)!=0x00) {
			MN.Right.SetLeaf(path, depth+1, check)		
		} else {
			MN.Left.SetLeaf(path, depth+1, check)
		}
	}

}

func (MN *MerkleNode) Instantiate(depth uint8) {
	
	if(depth>0){
		MN.Left=&MerkleNode{}
		MN.Left.Instantiate(depth-1)
		MN.Right=&MerkleNode{}
		MN.Right.Instantiate(depth-1)
	}
	return
}

func (MN *MerkleNode) CalcHash() (hash b512){
	if(MN.Left!=nil){
		MN.Data=Hash(MN.Left.CalcHash(),MN.Right.CalcHash())
		hash=MN.Data
		return
	}
	hash=MN.Data
	return
}
func (MT *MerkleTree) CalcHash() (hash b512) {
	MT.ValidHash=MT.Root.CalcHash()
	hash=MT.ValidHash
	return
}
func (MT *MerkleTree) ValidateHash(path uint8, check b512) (valid bool) {
	valid=b512eq(MT.Root.ValidateLeaf(path, 0, check),MT.ValidHash)
	return
}
func (MT *MerkleTree) Instantiate(depth uint8) {
	MT.Root=&MerkleNode{}
	MT.Root.Instantiate(depth)
}

func main(){
		fmt.Println("aaaaaaaaaaaaaaa")
}*/
