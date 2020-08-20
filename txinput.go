package simpleBlockchain

import (
	"bytes"
	"encoding/json"
)

var (
	coinbasePrevTxHash = []byte{
		0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
		0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
		0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
		0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
	}
	coinbasePrevTxOutIndex uint = 0
)

type TxIn struct {
	PrevTxHash Hashes		`json:"prevtxhash"`
	PrevTxOutIndex uint		`json:"prevtxoutindex"`
	ScriptSig Hashes		`json:"scriptsig"`
}

func CreateCoinbaseTxIn(data string) *TxIn {
	coinbaseTxIn := &TxIn{
		PrevTxHash: coinbasePrevTxHash,
		PrevTxOutIndex: coinbasePrevTxOutIndex,
		ScriptSig: []byte(data),
	}
	return coinbaseTxIn
}
func (txIn *TxIn) recoverCoinbaseScriptsig() string{
	return string(txIn.ScriptSig)
}

func (txIn *TxIn) isCoinbaseTxIn() bool {
	if bytes.Compare(txIn.PrevTxHash, coinbasePrevTxHash) == 0 && txIn.PrevTxOutIndex == coinbasePrevTxOutIndex {
		return true
	}
	return false
}

func (txIn *TxIn) Serialize() ([]byte, error) {
	res, err := json.Marshal(txIn)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DeserializeTxIn(data []byte) (*TxIn, error) {
	var txIn TxIn
	err := json.Unmarshal(data, &txIn)
	if err != nil {
		return nil, err
	}
	return &txIn, nil
}

