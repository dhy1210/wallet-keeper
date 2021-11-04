package api

import (
	"fmt"
	"net/http"

	"github.com/dhy1210/wallet-keeper/keeper"
	"github.com/dhy1210/wallet-keeper/keeper/btc"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (api *ApiServer) GetAddress(c *gin.Context) {
	value, _ := c.Get(KEEPER_KEY) // sure about the presence of this value
	keeper := value.(keeper.Keeper)

	account, found := c.GetQuery("account")
	if !found {
		account = btc.DEFAULT_ACCOUNT
	}
	address, err := keeper.GetAddress(account)
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, R(fmt.Sprint(err)))
	} else {
		c.JSON(http.StatusOK, R(address))
	}
}
