package simpleBlockchain

import (
	"math/big"
)

var GenesisBits = []byte{0x1f,0xff,0xff,0xff}
var GenesisTarget = []byte{
	0x00,0x00,0x00,0x00,0xff,0xff,0xff,0xff,0xff,0xff,
	0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
	0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,0xff,
	0xff,0xff,
}


type ProofOfWork struct {
	block *Block
	target *big.Int
}

func SerializeHeaderForMining(blockHeader *BlockHeader) []byte{
	return ConcatCopy(
		IntToLittleEndianBytes(blockHeader.Version),
		ReverseBytes(blockHeader.PrevBlock),
		ReverseBytes(blockHeader.MerkleRoot),
		IntToLittleEndianBytes(blockHeader.TimeStamp),
		IntToLittleEndianBytes(blockHeader.Bits),
		IntToLittleEndianBytes(blockHeader.Nonce),
	)
}


func NewProofOfWork(block *Block) *ProofOfWork{
	//bits := block.BlockHeader.Bits
	target := CalculateTarget(GenesisBits)
	//target := big.NewInt(0).SetBytes(GenesisTarget)
	return &ProofOfWork{
		block:  block,
		target: target,
	}
}

func (pow *ProofOfWork) mining(){
	var nonce = 0
	for {
		pow.block.BlockHeader.Nonce = uint32(nonce)
		res := big.NewInt(0).SetBytes(ReverseBytes(DoubleSha256(SerializeHeaderForMining(pow.block.BlockHeader))))
		if res.Cmp(pow.target) == -1 {
			return
		}
		nonce++
	}
}


// a = bits[0]
// b = bits[1]
// coefficient = bits[2:]
// target = coefficient * 2^(8*(a–3))
func CalculateTarget(bits []byte) *big.Int{
	//coefficient := big.NewInt(0).SetBytes(bits[2:])
	//a := big.NewInt(int64(bits[0]))
	//exp := big.NewInt(8).Mul(big.NewInt(8),a.Sub(a,big.NewInt(3)))
	//num := exp.Exp(big.NewInt(2), exp,nil)
	//target := coefficient.Mul(coefficient,num)
	//return target
	target := make([]byte, 32)
	exp := bits[0]
	coefficient := bits[1:]
	copy(target[32-exp:32-exp+3], coefficient)
	return big.NewInt(0).SetBytes(target)
}
