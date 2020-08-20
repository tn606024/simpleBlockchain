package simpleBlockchain

import "encoding/json"

const coinbaseReward = 5000000000

type TxOut struct {
	Value int			`json:"value"`
	ScriptPubKey Hashes `json:"scriptpubkey"`		// put public key hash
}

func CreateCoinbaseTxOut(address string) *TxOut{
	return &TxOut{
		Value:        coinbaseReward,
		ScriptPubKey: AddressToPubkeyHash(address),
	}
}

func (txOut *TxOut) isCoinbaseTxOut() bool {
	return txOut.Value == coinbaseReward
}

// Public_K=G Private_K=(x,y)
// Address=(Network Version) & Ripemd160(sha256(x&y) & checksum
// Checksum=First four bytes of sha256(sha256((Network Version)&Ripemd160(sha256(x&y))
// address base58((0x00||pubkeyHash||checksum(4bytes)))
func AddressToPubkeyHash(address string) []byte {
	decodeAddr := Base58Decode([]byte(address))
	return decodeAddr[1:len(decodeAddr)-4]
}

func (txOut *TxOut) Serialize() ([]byte, error) {
	res, err := json.Marshal(txOut)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DeserializeTxOut(data []byte) (*TxOut, error) {
	var txOut TxOut
	err := json.Unmarshal(data, &txOut)
	if err != nil {
		return nil, err
	}
	return &txOut, nil
}

