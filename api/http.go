package api

import (
	"fmt"
	"strings"

	"github.com/dhy1210/wallet-keeper/keeper"
	"github.com/dhy1210/wallet-keeper/keeper/btc"
	"github.com/dhy1210/wallet-keeper/keeper/eth"
	"github.com/dhy1210/wallet-keeper/keeper/usdt"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const KEEPER_KEY = "keeper"
const COIN_TYPE_HEADER = "CoinType"

// http api list
var METHODS_SUPPORTED = map[string]string{
	// misc
	"/ping":   "check if api service valid and backend bitcoin service healthy",
	"/health": "check system status",
	"/help":   "display this message",

	// useful APIs here
	"/getaddress":     "return address of specified account or default",
	"/getbalance":     "sum balances of all accounts",
	"/listaccounts":   "list accounts with amount, minconf is 6",
	"/getaccountinfo": "return account with corresponding bablance and addresses",
	"/createaccount":  "create account and return receive address, error if account exists",
	"/sendfrom":       "send amount of satoshi from some account to targets address",
	"/move":           "move from one account to another",
	// private
	//"/getnewaddress":  "return a new address of specified account or default",
	//"/getblockcount":  "return height of the blockchain",
	//"/listunspentmin": "list all unspent transactions",
	//"/sendtoaddress":  "Deprecicated: send amount of satoshi to address",
	//"/getaddress_with_balances": "all addresses together with balances",
}

type ApiServer struct {
	httpListenAddr string
	engine         *gin.Engine
	btcKeeper      keeper.Keeper
	usdtKeeper     keeper.Keeper
	ethKeeper      keeper.Keeper
}

//TODO valid host is valid
func (api *ApiServer) InitBtcClient(host, user, pass, logDir string) (err error) {
	api.btcKeeper, err = btc.NewClient(host, user, pass, logDir)
	return err
}

func (api *ApiServer) InitUsdtClient(host, user, pass, logDir string, propertyId int64) (err error) {
	api.usdtKeeper, err = usdt.NewClient(host, user, pass, logDir, propertyId)
	return err
}

func (api *ApiServer) InitEthClient(host, walletDir, password, accountPath, logDir string) (err error) {
	api.ethKeeper, err = eth.NewClient(host, walletDir, password, accountPath, logDir)
	return err
}

//Check
func (api *ApiServer) KeeperCheck() (err error) {
	if api.btcKeeper != nil {
		err = api.btcKeeper.Ping()
		if err != nil {
			err = errors.Wrap(err, "btc: ")
		}
	}

	if api.usdtKeeper != nil {
		err = api.usdtKeeper.Ping()
		if err != nil {
			err = errors.Wrap(err, "usdt: ")
		}
	}

	if api.ethKeeper != nil {
		err = api.ethKeeper.Ping()
		if err != nil {
			err = errors.Wrap(err, "eth: ")
		}
	}

	return err
}

func NewApiServer(addr string) *ApiServer {
	apiServer := &ApiServer{
		httpListenAddr: addr,
	}

	// build gin.Engine and register routers
	apiServer.buildEngine()

	return apiServer
}

func (api *ApiServer) buildEngine() {
	r := gin.Default()

	// with midlleware determine which currency is active
	// within this very session, should be either `btc` or  `usdt`.
	// if none of these present, abort this request and return caller
	// with bad request instantly.
	r.Use(func(c *gin.Context) {
		coin_type := strings.ToLower(c.Request.Header.Get(COIN_TYPE_HEADER))
		switch coin_type {
		case "btc":
			c.Set(KEEPER_KEY, api.btcKeeper)
			break
		case "usdt":
			c.Set(KEEPER_KEY, api.usdtKeeper)
			break
		case "eth":
			c.Set(KEEPER_KEY, api.ethKeeper)
			break
		default:
			c.JSON(400, gin.H{"message": "no coin type specified, should be btc or usdt"})
		}
	})

	// APIs related to wallet
	r.GET("/getblockcount", api.GetBlockCount)
	r.GET("/getaddress", api.GetAddress)
	r.GET("/getaddressesbyaccount", api.GetAddressesByAccount)
	r.GET("/getnewaddress", api.GetNewAddress)
	r.GET("/createaccount", api.CreateAccount)
	r.GET("/getaccountinfo", api.GetAccountInfo)
	r.GET("/listaccounts", api.ListAccounts)
	r.GET("/sendtoaddress", api.SendToAddress)
	r.GET("/sendfrom", api.SendFrom)
	r.GET("/listunspentmin", api.ListUnspentMin)
	r.GET("/move", api.Move)
	// r.GET("/test", api.Test)

	// misc API
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/health", func(c *gin.Context) {
		err := api.KeeperCheck()
		if err != nil {
			c.JSON(500, gin.H{
				"message": fmt.Sprint(err),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "healthy",
			})
		}
	})

	r.GET("/help", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"methods": METHODS_SUPPORTED,
		})
	})

	api.engine = r
}

func (api *ApiServer) HttpListen() error {
	return api.engine.Run(api.httpListenAddr)
}
