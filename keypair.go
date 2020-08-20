package simpleBlockchain

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"golang.org/x/crypto/ripemd160"
	"math/big"
)


const (
	networkVersion = byte(0x00)
	checksumLength = 4
)

type KeyPair struct {
	PublicKey []byte `json:"publickey"`
	PrivateKey *big.Int `json:"privatekey"`
}

func (keypair *KeyPair) getECDSAPublicKey() ecdsa.PublicKey{
	return ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:	big.NewInt(0).SetBytes(keypair.PublicKey[0:31]),
		Y:  big.NewInt(0).SetBytes(keypair.PublicKey[32:]),
	}
}

func byteToPublicKey(data []byte) *ecdsa.PublicKey{
	return &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:	big.NewInt(0).SetBytes(data[0:32]),
		Y:  big.NewInt(0).SetBytes(data[32:]),
	}
}

func (keypair *KeyPair) getECDSAPrivateKey() *ecdsa.PrivateKey{
	return &ecdsa.PrivateKey{
		PublicKey: keypair.getECDSAPublicKey(),
		D: keypair.PrivateKey,
	}
}

func (keypair *KeyPair) getAddress() (string, error){
	addr, err:= publicKeyToAddr(keypair.PublicKey)
	if err!= nil {
		return "", err
	}
	return addr, nil
}

func NewKeypair() *KeyPair {
	var keyPair KeyPair
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	pub := append(priv.PublicKey.X.Bytes(), priv.PublicKey.Y.Bytes()...)
	keyPair.PrivateKey = priv.D
	keyPair.PublicKey = pub
	return &keyPair
}


//https://en.bitcoin.it/wiki/File:PubKeyToAddr.png
func publicKeyToAddr(publicKey []byte) (string, error){
	publicKey = append([]byte{0x04}, publicKey...)
	shaPub := sha256.Sum256(publicKey)
	ripEncoder := ripemd160.New()
	_, err := ripEncoder.Write(shaPub[:])
	if err != nil {
		return "", err
	}
	hash := ripEncoder.Sum(nil)
	// append Network ID Byte
	hash = append([]byte{networkVersion}, hash...)
	//generate checksum
	shaHash := sha256.Sum256(hash)
	checkSum := sha256.Sum256(shaHash[:])
	address := Base58Encode(append(hash, checkSum[:checksumLength]...))
	return string(address), nil
}

func publicKeyToPublicKeyHash(publicKey []byte) string{
	addr, _ := publicKeyToAddr(publicKey)
	hash := AddressToPubkeyHash(addr)
	return hex.EncodeToString(hash)
}