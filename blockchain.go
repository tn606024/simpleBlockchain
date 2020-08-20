package simpleBlockchain

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"math/big"
)

var genesisBlock = &Block{
	BlockHeader: &BlockHeader{
		Version: 0,
		PrevBlock: Hashes{
			0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
			0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
			0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
			0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
		},
		MerkleRoot: Hashes{
			0xf8,0x20,0xe9,0x99,0xba,0x74,0x89,0x4e,
			0x26,0xd0,0xaf,0x26,0x1c,0xf2,0xa6,0xbc,
			0x4c,0xb1,0x19,0x0a,0x7b,0xff,0x4b,0xca,
			0x9b,0x32,0x83,0xad,0x3b,0x19,0xb2,0x02,
		},
		TimeStamp: 1597600039,
		Bits: 4294967071,
		Nonce: 89,
		Height: 1,
	},
	Transactions: []*Transaction{
		{
		Inputs: []*TxIn{
			{
				PrevTxHash:     Hashes{
					0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
					0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
					0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
					0x00,0x00,0x00,0x00,0x00,0x00,0x00,0x00,
				},
				PrevTxOutIndex: 0,
				ScriptSig:      Hashes{
					0x67,0x65,0x6e,0x65,0x73,0x69,0x73,0x20,
					0x62,0x6c,0x6f,0x63,0x6b,0x20,0x63,0x72,
					0x65,0x61,0x74,0x65,0x64,0x20,0x62,0x79,
					0x20,0x31,0x48,0x6d,0x70,0x34,0x44,0x44,
					0x4d,0x51,0x4b,0x39,0x75,0x78,0x51,0x6f,
					0x50,0x71,0x46,0x32,0x51,0x69,0x53,0x45,
					0x41,0x73,0x65,0x4e,0x33,0x66,0x6f,0x79,
					0x53,0x6f,0x77,
				},
			},
		},
		Outputs: []*TxOut{
			{
				Value:        5000000000,
				ScriptPubKey: Hashes{
					0xb7,0xfb,0x98,0x3d,0xac,0xeb,0x6f,0x73,
					0x16,0xa2,0x0b,0xb4,0xbf,0xc4,0x12,0xed,
					0x1f,0x5e,0x23,0x48,
				},
			},
		},
		LockTime: 0,
		},
	},
}



var dbSigName = "simpleBlockchain_%d.db"

type BlockChain struct {
	db *bolt.DB
	miner  string
	port   int
	height int
	top []byte
	isMining bool
	utxosMap map[string][]*UTXO
}

func NewBlockChain(address string, port int, isMining bool) *BlockChain {
	var bc *BlockChain
	exist := FindBlockchainExist(port)
	if exist == false {
		bc = CreateBlockChain(address, port, isMining)
		return bc
	}
	dbName := fmt.Sprintf(dbSigName,port)
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	bc = &BlockChain{
		db: db,
		miner: address,
		isMining: isMining,
	}
	err = db.View(func(tx *bolt.Tx) error {
		top := tx.Bucket([]byte("DB")).Get([]byte("top"))
		bc.top = top
		return nil
	})
	if err != nil {
		panic(err)
	}
	blk := bc.getBlockByHash(bc.top)
	bc.height = blk.BlockHeader.Height
	err = bc.ReIndexUTXO()
	if err != nil {
		panic(err)
	}
	utxos, err := bc.getUTXOs()
	if err != nil {
		panic(err)
	}
	bc.utxosMap = utxos
	return bc
}

func CreateBlockChain(address string, port int, isMining bool) *BlockChain {
	dbName := fmt.Sprintf(dbSigName, port)
	db, err := bolt.Open(dbName, 0600, nil)
	if err != nil {
		panic(err)
	}
	bc := &BlockChain{
		db: db,
		miner: address,
		isMining: isMining,
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte("DB"))
		if err != nil {
			panic(err)
		}
		_, err  = tx.CreateBucket([]byte("Index"))
		if err != nil {
			panic(err)
		}
		_, err = tx.CreateBucket([]byte("UTXO"))
		if err != nil {
			panic(err)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	//gblock := CreateGenesisBlock(address, fmt.Sprintf("genesis block created by %s", address))
	err = bc.AddBlock(genesisBlock)
	if err != nil {
		panic(err)
	}
	return bc
}

func (bc *BlockChain) ReOrg(blockHash []byte) error{
	err := bc.ReOrgBlockchain(blockHash)
	if err != nil {
		return err
	}
	err = bc.ReOrgUTXO()
	if err != nil {
		return err
	}
	return nil
}

func (bc *BlockChain) ReOrgBlockchain(blockHash []byte) error{
	hashes := bc.getBlockHashesAfterHash(blockHash)
	blk := bc.getBlockByHash(blockHash)
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB"))
		for _, hash := range hashes{
			err := b.Delete(hash)
			if err != nil {
				return err
			}
		}
		err := b.Put([]byte("top"), blockHash)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	bc.top = blockHash
	bc.height = blk.BlockHeader.Height
	return nil
}

func (bc *BlockChain) ReOrgUTXO() error{
	err := bc.db.Update(func(tx *bolt.Tx) error{
		butxo := tx.Bucket([]byte("UTXO"))
		err := butxo.Delete([]byte("top"))
		if err != nil{
			return err
		}
		return nil
	})
	bc.ReIndexUTXO()
	return err
}

func (bc *BlockChain) MiningEmptyBlock(miner string) (*Block, error){
	block := MiningNewBlock(miner, bc.top, GenesisBits, bc.height+1, []*Transaction{})
	err := bc.AddBlock(block)
	if err !=nil {
		return nil,err
	}
	return block, nil

}

func (bc *BlockChain) AddBlock(block *Block) error {
	if block.BlockHeader.Height == bc.height + 1 {
		err := bc.putBlock(block)
		if err != nil {
			return err
		}
		err = bc.ReIndexUTXO()
		if err != nil {
			return err
		}
		allutxos, err := bc.getUTXOs()
		if err != nil {
			return err
		}
		bc.utxosMap = allutxos
	}else {
		return fmt.Errorf("missing block height: %d", bc.height + 1)

	}
	return nil
}

func (bc *BlockChain) putBlock(block *Block) error{
	err := bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB"))
		bblock, err  := block.Serialize()
		if err != nil {
			return err
		}
		err = b.Put(block.newHash(), bblock)
		if err != nil {
			return err
		}
		if bc.height < block.BlockHeader.Height{
			bc.height = block.BlockHeader.Height
			bc.top = block.newHash()
			err := tx.Bucket([]byte("DB")).Put([]byte("top"), block.newHash())
			if err != nil{
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}


func (bc *BlockChain) ReIndexUTXO() error {
	var top []byte
	err := bc.db.View(func(tx *bolt.Tx) error {
		butxo := tx.Bucket([]byte("UTXO"))
		top = butxo.Get([]byte("top"))
		return nil
	})
	if err != nil {
		return err
	}
	if top == nil {
		utxosmap := bc.scanUTXOs()
		err := bc.db.Update(func(tx *bolt.Tx) error{
			butxo := tx.Bucket([]byte("UTXO"))
			err := butxo.Put([]byte("top"), bc.top)
			if err != nil{
				return err
			}
			for txid, utxos := range utxosmap{
				butxos, err  := json.Marshal(utxos)
				if err != nil {
					return err
				}
				butxo.Put(HexStrToBytes(txid),butxos)
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	}
	if bytes.Compare(top, bc.top) == 0 {
		fmt.Printf("UTXOs are already index to the best height\n")
		return nil
	}
	blks := bc.getBlockHashesAfterHash(top)
	for _, blk := range blks{
		block := bc.getBlockByHash(blk)
		for _, transaction  := range block.Transactions {
			var freshUtxos []*UTXO
			for _, input := range transaction.Inputs {
				if input.isCoinbaseTxIn() == false {
					err := bc.db.Update(func(tx *bolt.Tx) error {
						var utxos []*UTXO
						var newUtxos []*UTXO
						butxo := tx.Bucket([]byte("UTXO"))
						butxos := butxo.Get(input.PrevTxHash)
						err := json.Unmarshal(butxos, &utxos)
						if err != nil {
							return err
						}
						for _, utxo := range utxos {
							if utxo.Index != input.PrevTxOutIndex {
								newUtxos = append(newUtxos, utxo)
							}
						}
						bNewUtxos, err := json.Marshal(newUtxos)
						if err != nil {
							return err
						}
						butxo.Put(input.PrevTxHash, bNewUtxos)
						return nil
					})
					if err != nil {
						return err
					}
				}
			}
			for index, output := range transaction.Outputs {
				freshUtxos = append(freshUtxos, &UTXO{
					output,
					uint(index),
					hex.EncodeToString(transaction.newHash()),
				})
			}
			bfreshUtxos, err := json.Marshal(freshUtxos)
			if err != nil {
				return err
			}
			err = bc.db.Update(func(tx *bolt.Tx) error {
				butxo := tx.Bucket([]byte("UTXO"))
				butxo.Put(transaction.newHash(), bfreshUtxos)
				butxo.Put([]byte("top"),block.newHash())
				return nil
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (bc *BlockChain) getBlocks(desc bool) []*Block {
	var blks []*Block
	iter := bc.NewBlockIterator()
	for  iter.hasNext(){
		blk := iter.Next()
		if desc == true{
			blks = append(blks, blk)
		} else {
			blks = append([]*Block{blk}, blks...)
		}
	}
	return blks
}

func (bc *BlockChain) getBlockHashes(desc bool) [][]byte{
	var hashes [][]byte
	iter := bc.NewBlockIterator()
	for  iter.hasNext(){
		blk := iter.Next()
		if desc == true{
			hashes = append(hashes, blk.newHash())
		} else {
			hashes = append([][]byte{blk.newHash()}, hashes...)
		}
	}
	return hashes
}

func (bc *BlockChain) getBlockByHash(hash []byte) *Block{
	var blk *Block
	bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("DB"))
		blk,_ = DeserializeBlock(b.Get(hash))
		return nil
	})
	return blk
}


func (bc *BlockChain) getBlockHashesAfterHash(hash []byte) [][]byte {
	var blockhashes [][]byte
	iter := bc.NewBlockIterator()
	for iter.hasNext(){
		block := iter.Next()
		if bytes.Compare(block.newHash(), hash) != 0 {
			blockhashes = append([][]byte{block.newHash()},blockhashes...)
		}else{
			break
		}
	}
	return blockhashes
}

func (bc *BlockChain) getUTXOs() (map[string][]*UTXO, error){
	allUtxos := make(map[string][]*UTXO, 0)
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("UTXO"))
		c := b.Cursor()
		for txid, utxosb := c.First(); txid != nil; txid, utxosb = c.Next() {
			if bytes.Compare(txid,[]byte("top")) == 0{
				continue
			}
			var utxos []*UTXO
			err := json.Unmarshal(utxosb, &utxos)
			if err != nil {
				return err
			}
			allUtxos[hex.EncodeToString(txid)] = append(allUtxos[hex.EncodeToString(txid)], utxos...)
		}
		return nil
	})
	if err != nil{
		return nil, err
	}
	return allUtxos, nil
}

func (bc *BlockChain) findTransaction(searchtx []byte) *Transaction{
	iter := bc.NewBlockIterator()
	for iter.hasNext(){
		block := iter.Next()
		for _, tx := range block.Transactions {
			if bytes.Compare(tx.newHash(),searchtx)  == 0 {
				return tx
			}
		}
	}
	return nil
}

func (bc *BlockChain) scanUTXOs() map[string][]*UTXO{
	//var utxos []*UTXO
	//bc.db.View(func(tx *bolt.Tx) error {
	//	b := tx.Bucket([]byte("UTXO"))
	//	c := b.Cursor()
	//	for txid, utxoib := c.First(); txid != nil; txid, utxoib = c.Next() {
	//		utxoi := DeserializeUTXOIndex(utxoib)
	//		for _, unspent := range utxoi.UnspentOutputs{
	//			utxos = append(utxos, &UTXO{
	//				unspent: unspent,
	//				txid:    txid,
	//				height:  utxoi.Height,
	//			})
	//		}
	//	}
	//	return nil
	//})
	//return utxos
	iter := bc.NewBlockIterator()
	utxos := make(map[string][]*UTXO)
	spendUtxos := make(map[string][]int)
	for iter.hasNext(){
		block := iter.Next()
		for _, tx := range block.Transactions {
			txid :=  hex.EncodeToString(tx.newHash())
			loop:
			for index, out := range tx.Outputs{
				if spendUtxos[txid] != nil {
					ids := spendUtxos[txid]
					for _, id := range ids {
						if index == id {
							continue loop
						}
					}
				}
				utxos[txid] = append(utxos[txid],&UTXO{
					out,
					uint(index),
					txid,
				})
			}
			for index, in := range tx.Inputs {
				if tx.isCoinBase() == false{
					txinId := hex.EncodeToString(in.PrevTxHash)
					spendUtxos[txinId] = append(spendUtxos[txinId], index)
				}
			}
		}
	}
	return utxos
}

func (bc *BlockChain) verifyTransaction(transaction *Transaction) bool{
	if transaction.isCoinBase() == true {
		return true
	}
	for i, in:= range transaction.Inputs {
		signature := in.ScriptSig[:64]
		pubkey := in.ScriptSig[64:]
		isUnspent := false
		var prevTx *TxOut
		utxos := bc.utxosMap[hex.EncodeToString(in.PrevTxHash)]
		for _, unspent := range utxos {
			if unspent.Index == in.PrevTxOutIndex {
				isUnspent = true
				prevTx = unspent.Unspent
			}
		}
		if isUnspent == false {
			return false
		}
		message := transaction.CopyCleanScriptSigTx()
		message.Inputs[i].ScriptSig = prevTx.ScriptPubKey
		bmessage, _:= message.Serialize()
		r := big.NewInt(0).SetBytes(signature[:len(signature)/2])
		s := big.NewInt(0).SetBytes(signature[len(signature)/2:])
		isVerified := ecdsa.Verify(byteToPublicKey(pubkey),bmessage,r,s)
		if isVerified == false {
			return isVerified
		}
	}
	return true
}

func (bc *BlockChain) mining(miner string, bits []byte, transactions []*Transaction) (*Block, error) {
	block := MiningNewBlock(miner, bc.top, bits, bc.height+1, transactions)
	fmt.Println("Mining new block done")
	err := bc.AddBlock(block)
	fmt.Println("blockchain add new block now")
	if err != nil {
		return nil, err
	}
	return block, err
}

func FindBlockchainExist(port int) bool {
	dbName := fmt.Sprintf(dbSigName,port)
	exist := IsFileExists(dbName)
	return exist
}

func (bc *BlockChain) NewBlockIterator() *BlockIterator {
	var val *Block
	err := bc.db.View(func(tx *bolt.Tx) error {
		curr := tx.Bucket([]byte("DB")).Get(bc.top)
		currBlk,_ := DeserializeBlock(curr)
		val = currBlk
		return nil
	})
	if err != nil {
		panic(err)
	}
	iter := &BlockIterator{
		bc: bc,
		current: val,
	}
	return iter
}

type BlockIterator struct {
	bc *BlockChain
	current	*Block
}

func (bi *BlockIterator) hasNext() bool{
	if bi.current != nil {
		return true
	}
	return false
}

func (bi *BlockIterator) Next() *Block{
	var val *Block
	bi.bc.db.View(func(tx *bolt.Tx) error {
		var currBlk *Block
		curr := tx.Bucket([]byte("DB")).Get(bi.current.BlockHeader.PrevBlock)
		if curr != nil {
			currBlk,_ = DeserializeBlock(curr)
		}
		val = currBlk
		return nil
	})
	prev := bi.current
	bi.current = val
	return prev
}








