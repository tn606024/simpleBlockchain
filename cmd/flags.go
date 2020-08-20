package cmd

import (
	"github.com/urfave/cli/v2"
)

var (
	nodeportFlag = &cli.IntFlag{
		Name:	"nodeport",
		Usage:	"nodeport",
		Value:	3000,
	}
	apiportFlag = &cli.IntFlag{
		Name:	"apiport",
		Usage:	"apiport",
		Value:	8080,
	}
	minerFlag = &cli.StringFlag{
		Name:	"miner",
		Usage:  "miner address",
	}
	walletnameFlag = &cli.StringFlag{
		Name:	"walletname",
		Usage:  "wallet name",
		Required: true,
	}
	rawFlag = &cli.StringFlag{
		Name:	"raw",
		Usage: 	"need send raw",
	}
	isminingFlag = &cli.BoolFlag{
		Name:	"ismining",
		Usage:  "ismining",
		Required: true,
	}
	toFlag = &cli.StringFlag{
		Name:	"to",
		Usage:	"to address",
		Required: true,

	}
	amountFlag = &cli.IntFlag{
		Name:	"amount",
		Usage:	"send amount",
		Required: true,
	}

)