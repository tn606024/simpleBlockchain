package simpleBlockchain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Conn struct {
	conn	*http.Client
	url		string
}

type Response struct {
	Result json.RawMessage `json:"result"`
}

func NewConn(url string) *Conn{
	if url == "" {
		url = "http://127.0.0.1:8080"
	}
	return &Conn{
		conn: &http.Client{
			Timeout: 30 * time.Second,
		},
		url: url,
	}
}

func (c *Conn) get(route string, result interface{}) error {
	req, err := http.NewRequest("GET",fmt.Sprintf("%s/%s", c.url, route), nil)
	if err != nil{
		return fmt.Errorf("consturct http request error: %v\n", err)
	}
	res, err := c.conn.Do(req)
	if err != nil {
		return fmt.Errorf("connected error: %v\n", err)
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	switch res.StatusCode {
	case 200:
		var response Response
		err = json.Unmarshal(resBody, &response)
		if err != nil {
			return fmt.Errorf("json unmarshal resbody error: %v\n", err)
		}
		err = json.Unmarshal(response.Result, result)
		if err != nil {
			return fmt.Errorf("json Unmarshal response result error: %v\n", err)
		}
		return nil
	case 500:
		return fmt.Errorf("%s\n",string(resBody))
	}
	return nil
}

func (c *Conn) post(route string, result interface{}, msg interface{}) error {
	body, err := json.Marshal(msg)
	if err != nil{
		fmt.Errorf("json marshal error: %v", err)
	}
	req, err := http.NewRequest("POST",fmt.Sprintf("%s/%s", c.url, route), bytes.NewReader(body))
	req.Header.Set("Content-Type","application/json")
	if err != nil{
		return fmt.Errorf("consturct http request error: %v\n", err)
	}
	res, err := c.conn.Do(req)
	if err != nil {
		return fmt.Errorf("connected error: %v\n", err)
	}
	defer res.Body.Close()
	resBody, err := ioutil.ReadAll(res.Body)
	switch res.StatusCode {
	case 500:
		return fmt.Errorf("Server internal occured error: %s", string(resBody))
	case 400:
		return fmt.Errorf("Bad Request error: %s", string(resBody))
	}
	if err != nil {
		return err
	}
	var response Response
	err = json.Unmarshal(resBody, &response)
	if err != nil {
		return fmt.Errorf("json unmarshal resbody error: %v\n", err)
	}
	err = json.Unmarshal(response.Result, result)
	if err != nil {
		return fmt.Errorf("json Unmarshal response result error: %v\n", err)
	}
	return nil
}

func (c *Conn) MiningBlock() (block Block, err error){
	err = c.get("chain/mining", &block)
	return
}

func (c *Conn) GetBlocks() (blocks []*Block, err error) {
	err = c.get("chain/blocks", &blocks)
	return
}

func (c *Conn) GetBlockHashes() (hashes [][]byte, err error){
	err = c.get("chain/hashes", &hashes)
	return
}

func (c *Conn) GetBlockHeight() (height int, err error){
	err = c.get("chain/height", &height)
	return
}

func (c *Conn) GetUTXOs() (utxos []*UTXO, err error){
	err = c.get("chain/utxos", &utxos)
	return
}

func (c *Conn) GetWalletAddress() (addresses []string, err error){
	err = c.get("wallet/address", &addresses)
	return
}

func (c *Conn) GetWalletUTXOs() (utxos []*UTXO, err error){
	err = c.get("wallet/utxos", &utxos)
	return
}

func (c *Conn) GetWalletBalance() (balance int, err error){
	err = c.get("wallet/balance", &balance)
	return
}

func (c *Conn) SendTransaction(txobj TransactionObj) (tx Transaction, err error){
	err = c.post("wallet/send", &tx, txobj)
	return
}