package api

import (
	"fmt"
	"net/http"

	"github.com/dhy1210/wallet-keeper/keeper"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (api *ApiServer) GetBlockCount(c *gin.Context) {
	value, _ := c.Get(KEEPER_KEY) // sure about the presence of this value
	keeper := value.(keeper.Keeper)

	height, err := keeper.GetBlockCount()
	if err != nil {
		log.Error(err)
		c.JSON(http.StatusInternalServerError, R(fmt.Sprint(err)))
	} else {
		c.JSON(http.StatusOK, R(height))
	}
}
