package simpleBlockchain

import (
	"github.com/gin-gonic/gin"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"sync"
)

const (
	blockchainVersion = 1
	protocal = "tcp"
)

var KnownNodes = []string{
	"localhost:3000",
	"localhost:3001",
}

var knownNodeName = "knownnodes_%d.txt"

type MessageHeader string

const (
	VersionMsgHeader    MessageHeader =  "version"
	VerackMsgHeader 	  MessageHeader =  "verack"
	AddrMsgHeader   	  MessageHeader =  "addr"
	InvMsgHeader       MessageHeader =  "inv"
	GetDataMsgHeader	  MessageHeader	=  "getdata"
	GetBlocksMsgHeader  MessageHeader =  "getblocks"
	TxMsgHeader		  MessageHeader =  "tx"
	BlockMsgHeader  MessageHeader =  "block"
)

type logmsg interface {
	String() string
}

type Msg struct {
	Header			MessageHeader	`json:"header"`
	Payload			json.RawMessage	`json:"payload"`
}

type VersionMsg struct {
	Version 		int				`json:"version"`
	AddrFrom		string			`json:"addr_from"`
	StartHeight		int				`json:"start_height"`
}

func (msg *VersionMsg) String() string{
	bmsg, _ := json.MarshalIndent(msg,"","	")
	return string(bmsg) + "\n"
}

type VerackMsg struct {
	AddrFrom		string			`json:"addr_from"`
}

func (msg *VerackMsg) String() string{
	bmsg, _ := json.MarshalIndent(msg,"","	")
	return string(bmsg) + "\n"
}

type AddrMsg struct {
	AddrList		string			`json:"addr_list"`
}

func (msg *AddrMsg) String() string{
	bmsg, _ := json.MarshalIndent(msg,"","	")
	return string(bmsg) + "\n"
}

type InvMsg struct{
	AddrFrom	string			`json:"addr_from"`
	Type		string			`json:"type"`
	Hash		[]Hashes		`json:"hash"` // 32 byte
}

func (msg *InvMsg) String() string{
	bmsg, _ := json.MarshalIndent(msg,"","	")
	return string(bmsg) + "\n"
}

type GetdataMsg struct {
	AddrFrom	string			`json:"addr_from"`
	Type		string			`json:"type"`
	Hash		[]Hashes		`json:"hash"` // 32 byte
}

func (msg *GetdataMsg) String() string{
	bmsg, _ := json.MarshalIndent(msg,"","	")
	return string(bmsg) + "\n"
}

type GetblocksMsg struct {
	AddrFrom		string			`json:"addr_from"`
	BlockHashs		[]Hashes		`json:"blockhashes"`
}

func (msg *GetblocksMsg) String() string{
	bmsg, _ := json.MarshalIndent(msg,"","	")
	return string(bmsg) + "\n"
}

type TxMsg struct {
	AddrFrom	string				`json:"addr_from"`
	Transaction Hashes				`json:"transaction"`
}

func (msg *TxMsg) String() string{
	bmsg, _ := json.MarshalIndent(msg,"","	")
	return string(bmsg) + "\n"
}

type BlockMsg struct {
	AddrFrom	string					`json:"addr_from"`
	Block		[]Hashes				`json:"block"`
}

func (msg *BlockMsg) String() string{
	bmsg, _ := json.MarshalIndent(msg,"","	")
	return string(bmsg) + "\n"
}

type InvVect struct {
	Type		string	`json:"type"`
	Hash		Hashes	`json:"hash"` // 32 byte
}


type TransactionObj struct {
	To     string
	Amount int
}


type Server struct {
	node 		string
	wallet		*Wallet
	utxos		[]*UTXO
	nodeport    int
	apiport		int
	knownNodes	[]string
	connectMap	map[string]bool
	blockMap	map[string]int
	blockchain	*BlockChain
	mutex		sync.Mutex
}



func NewServer(nodeport int,  apiport int,  walletName string, isMining bool) *Server{
	node := fmt.Sprintf("localhost:%d",nodeport)
	wallet, err := GetExistWallet(walletName)
	addrs, _:=wallet.getAddresses()
	blockchain := NewBlockChain(addrs[0],nodeport, isMining)
	if err != nil {
		panic(err)
	}
	knownNodes, err := NewKnownNodes(nodeport)
	if err != nil {
		panic(err)
	}
	connectMap := make(map[string]bool,0)
	blockMap := make(map[string]int,0)
	utxos := make([]*UTXO,0)
	s:=  &Server{
		node: node,
		wallet: wallet,
		utxos: utxos,
		nodeport: nodeport,
		apiport: apiport,
		knownNodes: knownNodes,
		blockchain: blockchain,
		connectMap: connectMap,
		blockMap: blockMap,
	}
	s.ScanWalletUTXOs()
	return s
}

func NewKnownNodes(port int) ([]string, error){
	var knownNodes []string
	nodefile := fmt.Sprintf(knownNodeName,port)
	if IsFileExists(nodefile) == false {
		bkn, _ := json.Marshal(KnownNodes)
		err := ioutil.WriteFile(nodefile, bkn, 0644)
		if err != nil{
			return nil, err
		}
	}
	b, err := ioutil.ReadFile(fmt.Sprintf(knownNodeName,port))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &knownNodes)
	if err != nil {
		return nil, err
	}
	return knownNodes, nil
}

func (s *Server) StartServer() {
	ln, err := net.Listen(protocal, s.node)
	if err != nil {
		panic(err)
	}
	defer ln.Close()
	go s.StartApiServer(s.apiport)
	s.broadcastVersion()
	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}
		go s.handleConnection(conn)
	}
}

func (s *Server) StartApiServer(apiport int) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/chain/blocks", func(c *gin.Context){
		blks := s.blockchain.getBlocks(false)
		c.JSON(http.StatusOK, gin.H{
			"result": blks,
		})
	})
	r.GET("/chain/hashes", func(c *gin.Context){
		hashes := s.blockchain.getBlockHashes(false)
		c.JSON(http.StatusOK, gin.H{
			"result": hashes,
		})
	})
	r.GET("/chain/height", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"result": s.blockchain.height,
		})
	})
	r.GET("/chain/utxos", func(c *gin.Context){
		var allUtxos []*UTXO
		utxosMap, _  := s.blockchain.getUTXOs()
		for _, utxos := range utxosMap {
			for _, utxo := range utxos {
				allUtxos = append(allUtxos, utxo)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"result": allUtxos,
		})
	})
	r.GET("/chain/mining", func(c *gin.Context){
		blk, err := s.MiningEmptyBlockAndBroadcast()
		if err != nil {
			c.String(http.StatusInternalServerError, "server error occured: %s", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"result": blk,
		})
	})
	r.GET("/wallet/utxos", func(c *gin.Context){
		c.JSON(http.StatusOK, gin.H{
			"result": s.utxos,
		})
	})
	r.GET("/wallet/balance", func(c *gin.Context){
		balance := s.GetWalletBalance()
		c.JSON(http.StatusOK, gin.H{
			"result": balance,
		})
	})
	r.GET("/wallet/address", func(c *gin.Context){
		addrs, _ := s.wallet.getAddresses()
		c.JSON(http.StatusOK, gin.H{
			"result": addrs,
		})
	})
	r.POST("/wallet/send", func(c *gin.Context){
		var txObj TransactionObj
		c.BindJSON(&txObj)
		tx, err := s.SendTransaction(txObj.Amount, txObj.To)
		if err != nil {
			c.String(http.StatusInternalServerError, "server error occured: %s", err)
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"result": tx,
		})
	})
	r.Run(fmt.Sprintf(":%d",apiport))
}

func (s *Server) ScanWalletUTXOs() error{
	var utxos  []*UTXO
	allutxos, err := s.blockchain.getUTXOs()
	if err != nil {
		return err
	}
	pubkeys, _ := s.wallet.getPublickeyHash()
	for _, outs := range allutxos {
		for _, pubkey := range pubkeys {
			for _, out := range outs {
				if bytes.Compare(out.Unspent.ScriptPubKey, pubkey) == 0 {
					utxos= append(utxos, out)
				}
			}

		}
	}
	s.utxos = utxos
	return nil
}

func (s *Server) GetWalletBalance() int{
	sum := 0
	for _, utxo := range s.utxos {
		sum += utxo.Unspent.Value

	}
	return sum
}

func (s *Server) createTransaction(amount int, fee int, to string, change string) ([]byte, error){
	var uses []*UTXO
	var cost = 0
	s.ScanWalletUTXOs()
	for _, txout:= range s.utxos {
		if cost < amount+fee {
			uses = append(uses, txout)
			cost = cost + txout.Unspent.Value
		} else {
			break
		}
	}
	if cost < amount +fee {
		return nil, fmt.Errorf("you don't have enough coin")
	}
	tx := &Transaction{}
	for _, use := range uses {
		input := &TxIn{
			PrevTxHash:     HexStrToBytes(use.Txid),
			PrevTxOutIndex: use.Index,
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
	rawTx, err := s.wallet.signTransaction(tx)
	if err != nil {
		return nil, err
	}
	return rawTx, nil
}

func (s *Server) SendTransaction(amount int, to string) (*Transaction, error){
	btx, err := s.createTransaction(amount,0, to, s.blockchain.miner)
	if err != nil {
		return nil, err
	}
	tx, err := DeserializeTransaction(btx)
	if err != nil {
		return nil, err
	}
	isVerified := s.blockchain.verifyTransaction(tx)
	if isVerified == false {
		return nil, fmt.Errorf("transaction can't be verified")
	}
	if s.blockchain.isMining {
		block, err := s.blockchain.mining(s.blockchain.miner,GenesisBits,[]*Transaction{tx})
		if err != nil {
			return nil, err
		}
		s.broadcastBlock(block)
	}else {
		s.broadcastTx(tx)
	}
	return tx, nil
}

func (s *Server) MiningEmptyBlockAndBroadcast() (*Block,error) {
	if s.blockchain.isMining == false{
		return nil, fmt.Errorf("isMining is set false")
	}
	blk, err := s.blockchain.MiningEmptyBlock(s.blockchain.miner)
	if err != nil {
		return nil, err
	}
	s.broadcastBlock(blk)
	return blk, nil
}
func (s *Server) broadcastVersion(){
	for _, knownNode := range s.knownNodes {
		if knownNode != s.node {
			s.sendVersion(knownNode)
		}
	}
}

func (s *Server) broadcastBlock(block *Block){
	for _, knownNode := range s.knownNodes{
		if knownNode != s.node {
			s.sendBlock(knownNode, []*Block{block})
		}
	}
}

func (s *Server) broadcastTx(tx *Transaction){
	for _, knownNode := range s.knownNodes{
		if knownNode != s.node {
			s.sendTx(knownNode, tx)
		}
	}
}


func (s *Server) handleConnection(conn net.Conn){
	var msg Msg
	req, err := ioutil.ReadAll(conn)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(req, &msg)
	if err != nil{
		fmt.Printf("json unmarshal request error: %v\n", err)
	}
	switch msg.Header {
	case VersionMsgHeader:
		s.handleVersion(msg.Payload)
	case VerackMsgHeader:
		s.handleVerack(msg.Payload)
	case AddrMsgHeader:
	case InvMsgHeader:
		s.handleInv(msg.Payload)
	case GetDataMsgHeader:
		s.handleGetData(msg.Payload)
	case GetBlocksMsgHeader:
		s.handleGetBlocks(msg.Payload)
	case TxMsgHeader:
		s.handleTx(msg.Payload)
	case BlockMsgHeader:
		s.handleBlock(msg.Payload)
	}
}

func (s *Server) handleVersion(payload json.RawMessage){
	var versionMsg VersionMsg
	err := json.Unmarshal(payload, &versionMsg)
	if err != nil{
		fmt.Printf("json unmarshal error: %s\n", err)
	}
	logHandleMsg(VersionMsgHeader, &versionMsg)
	if versionMsg.Version == blockchainVersion {
		s.mutex.Lock()
		s.blockMap[versionMsg.AddrFrom] = versionMsg.StartHeight
		conn, ok := s.connectMap[versionMsg.AddrFrom]
		s.mutex.Unlock()
		if ok == true {
			if conn == false {
				s.sendVerack(versionMsg.AddrFrom)
			} else{
				s.sendVersion(versionMsg.AddrFrom)
			}
		} else {
			s.sendVersion(versionMsg.AddrFrom)
		}
	}
}

func (s *Server) handleVerack(payload json.RawMessage){
	var verackMsg VerackMsg
	err := json.Unmarshal(payload, &verackMsg)
	from := verackMsg.AddrFrom
	logHandleMsg(VerackMsgHeader, &verackMsg)
	if err != nil{
		fmt.Printf("json unmarshal error: %s\n", err)
	}
	s.mutex.Lock()
	height := s.blockchain.height
	conn, ok := s.connectMap[from]
	s.mutex.Unlock()
	if ok != true {
		fmt.Printf("this address doesn't exist in map: %s\n", from)
		return
	}
	if conn == false {
		s.sendVerack(from)
		return
	}
	s.mutex.Lock()
	block := s.blockMap[from]
	s.mutex.Unlock()
	if block > height {
		s.sendGetblocks(from)
	}
}

func (s *Server) handleGetBlocks(payload json.RawMessage){
	var getBlockMsg GetblocksMsg
	var invMsg	*InvMsg
	err := json.Unmarshal(payload, &getBlockMsg)
	if err != nil{
		fmt.Printf("json unmarshal error: %s\n", err)
	}
	logHandleMsg(GetBlocksMsgHeader, &getBlockMsg)
	hashes := getBlockMsg.BlockHashs
	myhashes := s.blockchain.getBlockHashes(false)
	reorgHeight := len(hashes) -1
	for i, hash := range hashes {
		if bytes.Compare(myhashes[i],(hash)) != 0 {
			if i == 0 {
				reorgHeight = 0
			}else {
				reorgHeight = i-1
			}
		}
	}


	//s.blockchain.ReOrg(reorgHash)
	invMsg = &InvMsg{
		Type: "block",
		Hash: bytesToHashes(myhashes)[reorgHeight:],
		AddrFrom: s.node,
	}
	s.sendInv(getBlockMsg.AddrFrom, invMsg)
}

func (s *Server) handleInv(payload json.RawMessage){
	var invMsg InvMsg
	var getDataMsg *GetdataMsg
	err := json.Unmarshal(payload, &invMsg)
	if err != nil {
		fmt.Printf("json unmarshal error: %s\n", err)
	}
	logHandleMsg(InvMsgHeader, &invMsg)
	if len(invMsg.Hash) > 1 {
		err := s.blockchain.ReOrg(invMsg.Hash[0])
		if err != nil {
			fmt.Printf("ReOrg blockchain happened error:%v", err)
		}
		fmt.Println("ReOrg done")
		switch invMsg.Type {
		case "block":
			getDataMsg = &GetdataMsg{
				AddrFrom: s.node,
				Type:     "block",
				Hash:     invMsg.Hash[1:],
			}
		case "tx":
			getDataMsg = &GetdataMsg{
				AddrFrom: s.node,
				Type:     "tx",
				Hash:     invMsg.Hash,
			}
		}
		s.sendGetData(invMsg.AddrFrom, getDataMsg)

	} else {
		switch invMsg.Type {
		case "block":
			getDataMsg = &GetdataMsg{
				AddrFrom: s.node,
				Type:     "block",
				Hash:     invMsg.Hash,
			}
		case "tx":
			getDataMsg = &GetdataMsg{
				AddrFrom: s.node,
				Type:     "tx",
				Hash:     invMsg.Hash,
			}
		}
		s.sendGetData(invMsg.AddrFrom, getDataMsg)
	}
}

func (s *Server) handleGetData(payload json.RawMessage){
	var getDataMsg GetdataMsg
	blocks := make([]*Block,0)
	err := json.Unmarshal(payload, &getDataMsg)
	if err != nil {
		fmt.Printf("json unmarshal error: %s\n", err)
	}
	logHandleMsg(GetDataMsgHeader, &getDataMsg)

	switch getDataMsg.Type {
	case "block":
		for _, hash := range getDataMsg.Hash {
			blk := s.blockchain.getBlockByHash(hash)
			blocks = append(blocks, blk)
		}
		s.sendBlock(getDataMsg.AddrFrom, blocks)

	case "tx":
		for _, hash := range getDataMsg.Hash {
			tx := s.blockchain.findTransaction(hash)
			s.sendTx(getDataMsg.AddrFrom, tx)
		}
	}
}


func (s *Server) handleBlock(payload json.RawMessage){
	var blockMsg BlockMsg
	blocks := make([]*Block,0 )
	err := json.Unmarshal(payload, &blockMsg)
	if err != nil {
		fmt.Printf("json unmarshal error: %s\n", err)
	}
	logHandleMsg(BlockMsgHeader, &blockMsg)
	for _, bblock := range blockMsg.Block {
		blk, _ := DeserializeBlock(bblock)
		blocks = append(blocks, blk)
	}
	if s.connectMap[blockMsg.AddrFrom] == true {
		for _, block := range blocks {
			err := s.blockchain.AddBlock(block)
			if err != nil {
				fmt.Println(err)
				return
			}
			s.ScanWalletUTXOs()
		}
	} else {
		s.sendVersion(blockMsg.AddrFrom)
	}
}

func (s *Server) handleTx(payload json.RawMessage) {
	var txMsg TxMsg
	var invMsg InvMsg
	err := json.Unmarshal(payload, &txMsg)
	if err != nil {
		fmt.Printf("json unmarshal error: %s\n", err)
	}
	logHandleMsg(TxMsgHeader, &txMsg)
	tx, _ := DeserializeTransaction(txMsg.Transaction)
	verify := s.blockchain.verifyTransaction(tx)
	if verify == false {
		fmt.Println("tx can't be verified")
		return
	}
	blk, err := s.blockchain.mining(s.blockchain.miner, GenesisBits, []*Transaction{tx})
	if err != nil {
		fmt.Printf("mining block occured error:%v", err)
	}
	invMsg = InvMsg{
		AddrFrom: s.node,
		Type: "block",
		Hash: []Hashes{blk.newHash()},
	}
	for _, knownNode := range s.knownNodes {
		if knownNode != s.node {
			s.sendInv(knownNode, &invMsg)
		}
	}

}

func (s *Server) sendVersion(addr string){
	versionMsg := VersionMsg{
		Version: 1,
		AddrFrom: s.node,
		StartHeight: s.blockchain.height,
	}
	data, err := contructMsg(VersionMsgHeader, versionMsg)
	if err != nil{
		fmt.Printf("%v", err)
		return
	}
	logSendMsg(VersionMsgHeader, addr, &versionMsg)
	err = s.send(addr, data)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	s.mutex.Lock()
	s.connectMap[addr] = false
	s.mutex.Unlock()
}

func (s *Server) sendVerack(addr string){
	verackMsg := VerackMsg{
		AddrFrom:s.node,
	}
	msg, err := contructMsg(VerackMsgHeader, verackMsg)
	if err != nil{
		fmt.Printf("%v", err)
		return
	}
	logSendMsg(VerackMsgHeader, addr, &verackMsg)
	err = s.send(addr, msg)
	if err == nil{
		s.mutex.Lock()
		s.connectMap[addr] = true
		s.mutex.Unlock()
	}
}

func (s *Server) sendGetblocks(addr string){
	s.mutex.Lock()
	hashes := s.blockchain.getBlockHashes(false)
	s.mutex.Unlock()
	getblocksMsg := GetblocksMsg{
		AddrFrom:   s.node,
		BlockHashs: bytesToHashes(hashes),
	}
	msg, err := contructMsg(GetBlocksMsgHeader, getblocksMsg)
	if err != nil{
		fmt.Printf("%v", err)
		return
	}
	logSendMsg(GetBlocksMsgHeader, addr, &getblocksMsg)
	s.send(addr, msg)
}

func (s *Server) sendInv(addr string, inv *InvMsg){
	msg, err := contructMsg(InvMsgHeader, inv)
	if err != nil{
		fmt.Printf("%v", err)
		return
	}
	logSendMsg(InvMsgHeader, addr, inv)
	s.send(addr, msg)
}

func (s *Server) sendGetData(addr string, getdatamsg *GetdataMsg){
	msg, err := contructMsg(GetDataMsgHeader, getdatamsg)
	if err != nil{
		fmt.Printf("%v", err)
		return
	}
	logSendMsg(GetDataMsgHeader, addr, getdatamsg)
	s.send(addr, msg)
}

func (s *Server) sendBlock(addr string, blks []*Block){
	bblks := make([]Hashes, 0)
	for _, blk := range blks {
		bblk, err := blk.Serialize()
		if err != nil{
			fmt.Printf("%v/n", err)
		}
		bblks = append(bblks, bblk)
	}

	blockMsg := BlockMsg{
		AddrFrom: s.node,
		Block: bblks,
	}
	msg, err := contructMsg(BlockMsgHeader, blockMsg)
	if err != nil{
		fmt.Printf("%v", err)
		return
	}
	logSendMsg(BlockMsgHeader, addr, &blockMsg)
	s.send(addr, msg)
}

func (s *Server) sendTx(addr string, tx *Transaction){
	btx, _:= tx.Serialize()
	txMsg := TxMsg{
		AddrFrom: s.node,
		Transaction: btx,
	}
	msg, err := contructMsg(TxMsgHeader, txMsg)
	if err != nil{
		fmt.Printf("%v", err)
		return
	}
	logSendMsg(TxMsgHeader,addr, &txMsg)
	s.send(addr, msg)
}

func (s *Server) send(addr string, data []byte) error{
	conn, err := net.Dial(protocal, addr)
	if err != nil {
		delete(s.blockMap, addr)
		delete(s.connectMap, addr)
		//s.deleteKnownNodes(addr)
		return fmt.Errorf("%s is not online \n", addr)
	}
	defer conn.Close()
	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		panic(err)
	}
	return nil
}

func contructMsg(header MessageHeader, payload interface{}) ([]byte, error){
	var msg Msg
	msg.Header = header
	if payload != nil {
		bpayload, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("json marshal error:%v", err)
		}
		msg.Payload = bpayload
	}
	data, err := json.Marshal(msg)
	if err != nil{
		return nil, fmt.Errorf("json marshal error:%v", err)
	}
	return data, nil
}

func logHandleMsg(header MessageHeader, msg logmsg){
	fmt.Printf("handle %s msg\n", header)
	fmt.Printf("msg:%s\n",msg.String())
}

func logSendMsg(header MessageHeader, to string, msg logmsg){
	fmt.Printf("send %s msg to %s\n", header, to)
	fmt.Printf("msg:%s\n",msg.String())
}