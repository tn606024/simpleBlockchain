package cmd

import (
	"fmt"
	"github.com/tn606024/simpleBlockchain"
	"github.com/urfave/cli/v2"
	"os"
)

var (
	startSubCommand = &cli.Command{
		Name:		 "start",
		Usage: 		 "start blockchain server",
		Description: "start blockchain server",
		ArgsUsage: 	 "<nodeport><apiport><walletname><ismining>",
		Flags: []cli.Flag{
			nodeportFlag,
			apiportFlag,
			walletnameFlag,
			isminingFlag,
		},
		Action: func(c *cli.Context) error {
			nodeport := c.Int("nodeport")
			apiport :=  c.Int("apiport")
			walletname := c.String("walletname")
			ismining := c.Bool("ismining")
			server := simpleBlockchain.NewServer(nodeport, apiport, walletname, ismining)
			server.StartServer()
			return nil
		},
	}
	miningblockSubCommand = &cli.Command{
		Name:		"miningblock",
		Usage:		"mine an empty block and broadcast to other node",
		Description: "mine an empty block and broadcast to other node",
		ArgsUsage: 	 "<apiport>",
		Flags: []cli.Flag{
			apiportFlag,
		},
		Action: func(c *cli.Context) error {
			apiport :=  c.Int("apiport")
			conn := simpleBlockchain.NewConn(fmt.Sprintf("http://127.0.0.1:%d", apiport))
			block, err := conn.MiningBlock()
			if err != nil {
				fmt.Printf("%v/n", err)
				os.Exit(1)
			}
			fmt.Println(block.String())
			return nil
		},
	}
	getblocksSubCommand = &cli.Command{
		Name:		"getblocks",
		Usage: 		 "get all blocks in blockchain",
		Description: "get all blocks in blockchain",
		ArgsUsage: 	 "<apiport>",
		Flags: []cli.Flag{
			apiportFlag,
		},
		Action: func(c *cli.Context) error {
			apiport :=  c.Int("apiport")
			conn := simpleBlockchain.NewConn(fmt.Sprintf("http://127.0.0.1:%d", apiport))
			blocks, err := conn.GetBlocks()
			if err != nil {
				fmt.Printf("%v/n", err)
				os.Exit(1)
			}
			for _, block := range blocks {
				fmt.Println(block.String())
			}
			return nil
		},
	}
	getblockhashesSubCommand = &cli.Command{
		Name:		"getblockhashes",
		Usage: 		 "get all blockhashes in blockchain",
		Description: "get all blockhashes in blockchain",
		ArgsUsage: 	 "<apiport>",
		Flags: []cli.Flag{
			apiportFlag,
		},
		Action: func(c *cli.Context) error {
			apiport :=  c.Int("apiport")
			conn := simpleBlockchain.NewConn(fmt.Sprintf("http://127.0.0.1:%d", apiport))
			hashes, err := conn.GetBlockHashes()
			if err != nil {
				fmt.Printf("%v/n", err)
				os.Exit(1)
			}
			for i, hash := range hashes {
				fmt.Printf("height: %d, hash:%s",i+1, simpleBlockchain.Hashes(hash).String())
			}
			return nil
		},
	}
	getblockheightSubCommand = &cli.Command{
		Name:		"getblockheight",
		Usage: 		 "get blockchain's height",
		Description: "get blockchain's height",
		ArgsUsage: 	 "<apiport>",
		Flags: []cli.Flag{
			apiportFlag,
		},
		Action: func(c *cli.Context) error {
			apiport :=  c.Int("apiport")
			conn := simpleBlockchain.NewConn(fmt.Sprintf("http://127.0.0.1:%d", apiport))
			height, err := conn.GetBlockHeight()
			if err != nil {
				fmt.Printf("%v/n", err)
				os.Exit(1)
			}
			fmt.Println(height)

			return nil
		},
	}

	getutxosSubCommand = &cli.Command{
		Name:		"getutxos",
		Usage: 		 "get all utxos in blockchain",
		Description: "get all utxos in blockchain",
		ArgsUsage: 	 "<apiport>",
		Flags: []cli.Flag{
			apiportFlag,
		},
		Action: func(c *cli.Context) error {
			apiport :=  c.Int("apiport")
			conn := simpleBlockchain.NewConn(fmt.Sprintf("http://127.0.0.1:%d", apiport))
			utxos, err := conn.GetUTXOs()
			if err != nil {
				fmt.Printf("%v/n", err)
				os.Exit(1)
			}
			for _, utxo := range utxos {
				fmt.Println(utxo.String())
			}
			return nil
		},
	}
	getwalletaddressSubCommand  = &cli.Command{
		Name:		"getwalletaddress",
		Usage: 		 "get wallet address set in server",
		Description: "get wallet address set in server",
		ArgsUsage: 	 "<apiport>",
		Flags: []cli.Flag{
			apiportFlag,
		},
		Action: func(c *cli.Context) error {
			apiport :=  c.Int("apiport")
			conn := simpleBlockchain.NewConn(fmt.Sprintf("http://127.0.0.1:%d", apiport))
			addrs, err := conn.GetWalletAddress()
			if err != nil {
				fmt.Printf("%v/n", err)
				os.Exit(1)
			}
			for _, addr := range addrs {
				fmt.Println(addr)
			}
			return nil
		},
	}
	getwalletutxosSubCommand  = &cli.Command{
		Name:		"getwalletutxos",
		Usage: 		 "get wallet's utxos",
		Description: "get wallet's utxos",
		ArgsUsage: 	 "<apiport>",
		Flags: []cli.Flag{
			apiportFlag,
		},
		Action: func(c *cli.Context) error {
			apiport :=  c.Int("apiport")
			conn := simpleBlockchain.NewConn(fmt.Sprintf("http://127.0.0.1:%d", apiport))
			utxos, err := conn.GetWalletUTXOs()
			if err != nil {
				fmt.Printf("%v/n", err)
				os.Exit(1)
			}
			for _, utxo := range utxos {
				fmt.Println(utxo.String())
			}
			return nil
		},
	}
	getwalletbalanceSubCommand  = &cli.Command{
		Name:		"getwalletbalance",
		Usage: 		 "get wallet's balance",
		Description: "get wallet's balance",
		ArgsUsage: 	 "<apiport>",
		Flags: []cli.Flag{
			apiportFlag,
		},
		Action: func(c *cli.Context) error {
			apiport :=  c.Int("apiport")
			conn := simpleBlockchain.NewConn(fmt.Sprintf("http://127.0.0.1:%d", apiport))
			balance, err := conn.GetWalletBalance()
			if err != nil {
				fmt.Printf("%v/n", err)
				os.Exit(1)
			}
			fmt.Println(balance)
			return nil
		},
	}
	sendTransactionSubCommand = &cli.Command{
		Name:	"sendTransaction",
		Usage:	"create Transaction and broadcast to other node",
		Description: "create Transaction and broadcast to other node",
		ArgsUsage: 	 "<apiport><to><amount>",
		Flags: []cli.Flag{
			apiportFlag,
			toFlag,
			amountFlag,
		},
		Action: func(c *cli.Context) error {
			apiport :=  c.Int("apiport")
			to := c.String("to")
			amount :=  c.Int("amount")
			conn := simpleBlockchain.NewConn(fmt.Sprintf("http://127.0.0.1:%d", apiport))
			tx, err := conn.SendTransaction(simpleBlockchain.TransactionObj{
				To:     to,
				Amount: amount,
			})
			if err != nil {
				fmt.Printf("%v/n", err)
				os.Exit(1)
			}
			fmt.Println(tx.String())
			return nil
		},
	}
	ServerCommand = &cli.Command{
		Name:	"server",
		Usage:	"blockchain server commands",
		ArgsUsage: "",
		Category: "Server Commands",
		Description: "",
		Subcommands: []*cli.Command{
			startSubCommand,
			getblocksSubCommand,
			getblockhashesSubCommand,
			getblockheightSubCommand,
			getutxosSubCommand,
			getwalletaddressSubCommand,
			getwalletutxosSubCommand,
			getwalletbalanceSubCommand,
			sendTransactionSubCommand,
			miningblockSubCommand,
		},
	}
)