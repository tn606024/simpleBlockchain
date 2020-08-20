package simpleBlockchain

import "encoding/json"

type BlockHeader struct {
	Version	 		uint32	`json:"version"`
	PrevBlock 		Hashes 	`json:"prevblock"`
	MerkleRoot 		Hashes	`json:"merkleroot"`
	TimeStamp 		uint32	`json:"timestamp"`
	Bits 			uint32	`json:"bits"`
	Nonce 			uint32	`json:"nonce"`
	Height			int		`json:"height"`
}



func (bh *BlockHeader) Serialize() ([]byte, error) {
	bbh, err := json.Marshal(bh)
	if err != nil {
		return nil, err
	}
	return bbh, nil
}

func copyBlockHeader(blockheader *BlockHeader) BlockHeader{
	newPrevBlock := make([]byte, len(blockheader.PrevBlock))
	newMerkleRoot := make([]byte, len(blockheader.MerkleRoot) )
	copy(newPrevBlock, blockheader.PrevBlock )
	copy(newMerkleRoot, blockheader.MerkleRoot)
	return BlockHeader{
		blockheader.Version,
		newPrevBlock,
		newMerkleRoot,
		blockheader.TimeStamp,
		blockheader.Bits,
		blockheader.Nonce,
		blockheader.Height,
	}

}

func DeserializeBlockHeader(data []byte) (*BlockHeader, error) {
	var bh BlockHeader
	err := json.Unmarshal(data, &bh)
	if err != nil {
		return nil, err
	}
	return &bh, nil
}