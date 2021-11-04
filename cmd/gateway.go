package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dhy1210/wallet-keeper/api"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type Backends string

func (b Backends) Contains(target string) bool {
	return strings.Contains(strings.ToLower(string(b)), strings.ToLower(target))
}

var gateCmd = cli.Command{
	Name:    "run",
	Aliases: []string{"r"},
	Flags: []cli.Flag{
		httpAddrFlag,
		btcRpcAddrFlag,
		btcRpcUserFlag,
		btcRpcPassFlag,
		usdtRpcAddrFlag,
		usdtRpcUserFlag,
		usdtRpcPassFlag,
		usdtPropertyIdFlag,
		ethRpcAddrFlag,
		ethWalletDirFlag,
		ethAccountPasswordFlag,
		ethAccountFlag,
		backendsFlag,
	},
	Usage: "serve api gateway",
	Action: func(c *cli.Context) error {

		var err error
		apiServer := api.NewApiServer(c.String("http-listen-addr"))

		backends := Backends(c.String("backends"))
		if backends.Contains("btc") {
			log.Infof("connecting to btc addr: %s", c.String("btc-rpc-addr"))
			err = apiServer.InitBtcClient(
				c.String("btc-rpc-addr"),  // host
				c.String("btc-rpc-user"),  // user
				c.String("btc-rpc-pass"),  // password
				c.GlobalString("log-dir"), // logDir
			)
			if err != nil {
				log.Error(err)
				return err
			}
		}

		if backends.Contains("usdt") {
			log.Infof("connecting to usdt addr: %s", c.String("usdt-rpc-addr"))
			err = apiServer.InitUsdtClient(
				c.String("usdt-rpc-addr"),        // host
				c.String("usdt-rpc-user"),        // user
				c.String("usdt-rpc-pass"),        // password
				c.GlobalString("log-dir"),        // logDir
				int64(c.Int("usdt-property-id")), // propertyId
			)
			if err != nil {
				log.Error(err)
				return err
			}
		}

		if backends.Contains("eth") {
			log.Infof("connecting to eth  addr: %s", c.String("eth-rpc-addr"))
			err = apiServer.InitEthClient(
				c.String("eth-rpc-addr"), // host
				c.String("eth-wallet-dir"),
				c.String("eth-account-password"),
				c.String("eth-account-path"),
				c.GlobalString("log-dir"), // logDir
			)
			if err != nil {
				log.Error(err)
				return err
			}
		}

		// Check btc/usdt rpc call connectivity
		err = apiServer.KeeperCheck()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return err
		}

		fmt.Fprintf(os.Stdout, "starting api gateway with addr: %s", c.String("http-listen-addr"))
		// start accepting http requests
		return apiServer.HttpListen()
	},
}
