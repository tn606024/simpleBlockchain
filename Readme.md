Simple Blockchain
------

This repository contains the golang code of a simple blockchain implementation.
  
This blockchain consists of three parts:  
- A simple wallet that you can get address, scan utxos, sign transaction.  
- A simple blockchain can sync block from other known nodes, mining new block, send transaction and broadcast to other node.  
- A simple restful server you can query blocks and utxos from blockchain.  

There are many part are not like real blockchain because it's just simple implementation, still
insecure and incomplete. you can learn the basic operation of the blockchain through this project.


How to run
------

## Build

```shell script
go build ./cmd/cli
```

### Create Wallet

```shell script
./cli wallet create -walletname "alice"
./cli wallet create -walletname "bob"
```

### Start two node

```shell script
./cli server start -nodeport 3000 -apiport 8080 -walletname "alice" -ismining=true
./cli server start -nodeport 3001 -apiport 8081 -walletname "bob" -ismining=true
```

### Mining empty block to get block reward
```shell script
./cli server miningblock --apiport 8080
```

### Send Transaction to other address
```shell script
./cli server sendTransaction --apiport 8080 --to "172wJyiJZxXWyBW7CYSVddsR5e7ZMxtja9" -amount 100000
```

threr are still have other blockchain command, you can find out by type `./cli server`.


Example
------

### Create wallet

```shell script
./cli wallet create -walletname "alice"
```

### Start blockchain server

 ```shell script
 ./cli server start -nodeport 3000 -apiport 8080 -walletname "alice" -ismining=true
 ```

### Get blocks

 ```shell script
 ./cli server getblocks -apiport 8080
 ```

### Get block hashes

 ```shell script
 ./cli server getblockhashes -apiport 8080
 ```

### Get block height

 ```shell script
 ./cli server getblockheight -apiport 8080
 ```

### Get block utxos

 ```shell script
 ./cli server getutxos -apiport 8080
 ```

### Get wallet address

 ```shell script
 ./cli server getwalletaddress -apiport 8080
 ```

### Get wallet utxos

 ```shell script
 ./cli server getwalletutxos -apiport 8080
 ```

### Get wallet balance

 ```shell script
 ./cli server getwalletbalance -apiport 8080
 ```

### Send transaction

 ```shell script
./cli server sendTransaction --apiport 8080 --to "172wJyiJZxXWyBW7CYSVddsR5e7ZMxtja9" -amount 100000
 ```

### Mining block

 ```shell script
 ./cli server miningblock -apiport 8080
 ```


