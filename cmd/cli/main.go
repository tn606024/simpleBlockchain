package main

import (
	"fmt"
	"github.com/ian/simpleBlockchain/cmd"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			cmd.ServerCommand,
			cmd.WalletCommand,
		},
	}


	err := app.Run(os.Args)
	if err != nil {
		fmt.Errorf("%s",err)
	}
}