package simpleBlockchain

import "encoding/json"

// https://btcinformation.org/en/developer-reference#raw-transaction-format
type Transaction struct {
	Inputs []*TxIn		`json:"inputs"`
	Outputs []*TxOut	`json:"outputs"`
	LockTime uint32 	`json:"locktime"`// 4 bytes
}

func (tx *Transaction) isCoinBase() bool{
	return len(tx.Inputs) == 1 && tx.Inputs[0].isCoinbaseTxIn() == true && len(tx.Outputs) == 1 &&  tx.Outputs[0].isCoinbaseTxOut() == true
}

func CreateCoinBaseTransaction(address string, data string) *Transaction{
	txIn := CreateCoinbaseTxIn(data)
	txOut := CreateCoinbaseTxOut(address)
	tx := &Transaction{
		Inputs:  []*TxIn{txIn},
		Outputs: []*TxOut{txOut},
	}
	return tx
}

func (tx *Transaction)CopyCleanScriptSigTx() *Transaction{
	var inputs  []*TxIn
	var outputs []*TxOut
	for _, input := range tx.Inputs{
		inputs = append(inputs, &TxIn{
			PrevTxHash:     input.PrevTxHash,
			PrevTxOutIndex: input.PrevTxOutIndex,
			ScriptSig:      nil,
		})
	}
	for _, output := range tx.Outputs{
		outputs = append(outputs, &TxOut{
			Value:        output.Value,
			ScriptPubKey: output.ScriptPubKey,
		})
	}
	return &Transaction{
		Inputs:   inputs,
		Outputs:  outputs,
		LockTime: tx.LockTime,
	}
}


func (tx *Transaction) newHash() []byte {
	txBytes,_ := tx.Serialize()
	return ReverseBytes(DoubleSha256(txBytes))
}

func (tx *Transaction) Serialize() ([]byte, error) {
	res, err := json.Marshal(tx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (tx *Transaction) String() string{
	btx, _ := json.MarshalIndent(tx,"","	")
	return string(btx) + "\n"
}

func DeserializeTransaction(data []byte) (*Transaction, error){
	var tx Transaction
	err := json.Unmarshal(data, &tx)
	if err != nil {
		return nil, err
	}
	return &tx, nil
}