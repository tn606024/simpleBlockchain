package simpleBlockchain

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"
)

var (
	coinbasePrevBlock = []byte{
		0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
		0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
		0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
		0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
	}
	genesisBlockPrevBlock= []byte{
		0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
		0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
		0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
		0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
	}
)

type Block struct {
	BlockHeader *BlockHeader	`json:"blockheader""`
	Transactions []*Transaction	`json:"transactions"`
}

func (b Block) newHash() []byte {
	bh, _ := b.BlockHeader.Serialize()
	return ReverseBytes(DoubleSha256(bh))
}

func CreateGenesisBlock(miner string, data string) *Block {
	return MiningNewBlock(miner, genesisBlockPrevBlock, GenesisBits,1, []*Transaction{})
}


func MiningNewBlock(miner string, prevBlock []byte, bits []byte, height int, transactions []*Transaction) *Block{
	coinbaseTx := CreateCoinBaseTransaction(miner, fmt.Sprintf("mine by %s at height %d",miner, height))
	txs := append([]*Transaction{coinbaseTx}, transactions...)
	root := CalculateMerkleRoot(txs)
	bh := &BlockHeader{
		Version:	0,
		PrevBlock: prevBlock,
		MerkleRoot: root,
		TimeStamp: uint32(time.Now().Unix()),
		Bits: binary.LittleEndian.Uint32(bits),
		Height: height,
	}
	block := &Block{
		BlockHeader:  bh,
		Transactions: txs,
	}
	pow := NewProofOfWork(block)
	pow.mining()
	newBlockHeader:= copyBlockHeader(pow.block.BlockHeader)
	pow.block.BlockHeader = &newBlockHeader
	return pow.block
}



func (block Block) Serialize() ([]byte,error) {
	bblock, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}
	return bblock, nil
}

func DeserializeBlock(data []byte) (*Block,error) {
	var blk Block
	err := json.Unmarshal(data, &blk)
	if err != nil {
		return nil, err
	}
	return &blk, nil
}

func (b Block) String() string {
	bs, _ := json.MarshalIndent(b,"","	")
	return string(bs) + "\n"
}