package simpleBlockchain

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

const (
	walletFile = "wallet_%s"
)

type Wallet struct {
	KeyPairs map[string]*KeyPair
}


func NewWallet(name string) (*Wallet, error){
	filename := fmt.Sprintf(walletFile, name)
	if IsFileExists(filename) == true {
		return nil, errors.New("this file exist!!")
	}
	keypair := NewKeypair()
	keyPairs := make(map[string]*KeyPair, 0)
	keyPairs[publicKeyToPublicKeyHash(keypair.PublicKey)] = keypair
	wallet := Wallet{
		KeyPairs:keyPairs,
	}
	jsonBytes, err := json.Marshal(wallet)
	if err != nil{
		return nil, err
	}
	err = ioutil.WriteFile(filename, jsonBytes, 0644)
	if err != nil {
		return nil, err
	}
	return &wallet, nil
}

func GetExistWallet(name string) (*Wallet, error){
	var wallet Wallet
	filename := fmt.Sprintf(walletFile, name)
	if IsFileExists(filename) == false {
		return nil, errors.New("this file doesn't exist!!")
	}
	jsonBytes, err:= ioutil.ReadFile(filename)
	if err!= nil {
		return nil, err
	}

	err = json.Unmarshal(jsonBytes, &wallet)
	if err!= nil {
		return nil, err
	}
	return &wallet, nil
}

func (wallet *Wallet) getAddresses() ([]string, error) {
	var addresses []string
	for _, keypair := range wallet.KeyPairs {
		addr, err := keypair.getAddress()
		if err!= nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}
	return addresses, nil
}

func (wallet *Wallet) getPublickeyHash() ([][]byte, error) {
	var hashes [][]byte
	addresses, err := wallet.getAddresses()
	if err!= nil {
		return nil, err
	}
	for _, addr := range addresses {
		hash := AddressToPubkeyHash(addr)
		hashes = append(hashes, hash)
	}
	return hashes, nil
}

func (wallet *Wallet) getPublickeys() [][]byte {
	var publicKeys [][]byte
	for _, keypair := range wallet.KeyPairs {
		pubkey := keypair.PublicKey
		publicKeys = append(publicKeys, pubkey)
	}
	return  publicKeys
}

func (wallet *Wallet) signMessageByKey(key string, message []byte) ([]byte, error){
	priv := wallet.KeyPairs[key].getECDSAPrivateKey()
	r, s, err := ecdsa.Sign(rand.Reader, priv, message)
	if err != nil {
		return nil, err
	}
	signature := append(r.Bytes(),s.Bytes()...)
	return signature, nil
}

func (wallet *Wallet) signTransaction(tx *Transaction) ([]byte, error) {
	var scriptSigs [][]byte
	scriptSigs = make([][]byte, len(tx.Inputs))
	for i, input := range tx.Inputs {
		key := hex.EncodeToString(input.ScriptSig)
		temp := input.ScriptSig
		cleanTx := tx.CopyCleanScriptSigTx()
		cleanTx.Inputs[i].ScriptSig = temp
		bcleanTx, _ := cleanTx.Serialize()
		signature, err := wallet.signMessageByKey(key,bcleanTx)
		if err != nil {
			return nil, err
		}
		scriptSigs[i] = append(signature, wallet.KeyPairs[key].PublicKey...)
	}
	for i, scriptSig := range scriptSigs {
		tx.Inputs[i].ScriptSig = scriptSig
	}
	btx, _:= tx.Serialize()
	return btx, nil
}


type UTXO struct {
	Unspent 	*TxOut		`json:"unspent"`
	Index		uint		`json:"index"`
	Txid		string		`json:"txid"`
}

func (utxo *UTXO) Serialize() ([]byte, error) {
	butxo, err := json.Marshal(utxo)
	if err != nil{
		return nil, err
	}
	return butxo, nil
}

func DeserializeUTXO(butxo []byte) (*UTXO, error){
	var utxo UTXO
	err := json.Unmarshal(butxo, &utxo)
	if err != nil {
		return nil, err
	}
	return &utxo, nil
}

func (utxo UTXO) String() string {
	butxo, _ := json.MarshalIndent(utxo,"","	")
	return string(butxo) + "\n"
}

type BlockchainWallet struct {
	wallet *Wallet
	server *Server
	utxos  map[string][]*UTXO

}

func NewBlockchainWallet(server *Server, wallet *Wallet) *BlockchainWallet{
	bw := &BlockchainWallet{
		wallet: wallet,
		server: server,
	}
	bw.ScanUTXOs()
	return bw

}

func (bw BlockchainWallet) GetBalance() int {
	sum := 0
	bw.ScanUTXOs()
	for _, utxos := range bw.utxos {
		for _, utxo := range utxos {
			sum += utxo.Unspent.Value
		}
	}
	return sum
}

func (bw BlockchainWallet) ScanUTXOs() error{
	var utxos  map[string][]*UTXO
	allutxos, err := bw.server.blockchain.getUTXOs()
	if err != nil {
		return err
	}
	pubkeys := bw.wallet.getPublickeys()
	for txid, outs := range allutxos {
		for _, pubkey := range pubkeys {
			for _, out := range outs {
				if bytes.Compare(out.Unspent.ScriptPubKey, pubkey) == 0 {
					utxos[txid] = append(utxos[txid], out)
				}
			}

		}
	}
	bw.utxos = utxos
	return nil
}

func (bw *BlockchainWallet) CreateTransaction(amount int, fee int, to string, change string) ([]byte, error){
	var uses []*UTXO
	var cost = 0
	if len(bw.utxos) == 0 {
		bw.ScanUTXOs()
	}
	for _, txouts:= range bw.utxos {
		for _, txout := range txouts {
			if cost < amount+fee {
				uses = append(uses, txout)
				cost = cost + txout.Unspent.Value
			} else {
				break
			}
		}
	}
	if cost < amount +fee {
		return nil, fmt.Errorf("you don't have enough coin")
	}
	var tx *Transaction
	for _, use := range uses {
		input := &TxIn{
			PrevTxHash:     HexStrToBytes(use.Txid),
			PrevTxOutIndex: uint(use.Index),
			ScriptSig:      use.Unspent.ScriptPubKey,
		}
		tx.Inputs = append(tx.Inputs, input)
	}
	out := &TxOut{
		Value:        amount,
		ScriptPubKey: AddressToPubkeyHash(to),
	}
	tx.Outputs = append(tx.Outputs, out)
	if cost - amount - fee > 0 {
		receive := &TxOut{
			Value: cost- amount - fee,
			ScriptPubKey: AddressToPubkeyHash(change),
		}
		tx.Outputs = append(tx.Outputs, receive)
	}

	rawTx, err := bw.wallet.signTransaction(tx)
	if err != nil {
		return nil, err
	}
	return rawTx, nil
}

func (bw *BlockchainWallet) Send(){

}



