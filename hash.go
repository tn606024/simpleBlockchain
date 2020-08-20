package simpleBlockchain

import (
	"encoding/hex"
	"encoding/json"
)

type Hashes []byte

func (h Hashes) MarshalJSON() ([]byte, error) {
	res, err:= json.Marshal(hex.EncodeToString(h))
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (h *Hashes) UnmarshalText(text []byte) (err error) {
	b, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	*h = b
	return
}

func (h Hashes) String() string {
	bs, _ := json.MarshalIndent(h,"","	")
	return string(bs) + "\n"
}

func HashestoBytes(h []Hashes) ([][]byte) {
	ans := make([][]byte, len(h))
	for i, hash := range h {
		ans[i] = hash
	}
	return ans
}

func bytesToHashes(inputs [][]byte) []Hashes{
	ans := make([]Hashes, len(inputs))
	for i, input := range inputs {
		ans[i] = input
	}
	return ans
}
