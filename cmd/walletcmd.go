package cmd

import (
	"fmt"
	"github.com/ian/simpleBlockchain"
	"github.com/urfave/cli/v2"
	"os"
)

var (
	createSubCommand = &cli.Command{
		Name:		 "create",
		Usage: 		 "create new wallet",
		Description: "create new wallet",
		ArgsUsage: 	 "<walletname>",
		Flags: []cli.Flag{
			walletnameFlag,
		},
		Action: func(c *cli.Context) error {
			walletname := c.String("walletname")
			_, err := simpleBlockchain.NewWallet(walletname)
			if err != nil {
				fmt.Printf("wallet create error:%v/n", err)
				os.Exit(1)
			}
			fmt.Println("create wallet success")

			return nil
		},
	}
	WalletCommand = &cli.Command{
		Name:	"wallet",
		Usage:	"wallet command",
		ArgsUsage: "",
		Category: "Wallet Commands",
		Description: "",
		Subcommands: []*cli.Command{
			createSubCommand,
		},
	}
)
