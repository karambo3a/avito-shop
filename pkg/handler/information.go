package handler

import (
	"avito_go/pkg/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) GetInformation(c *gin.Context) {
	id, err := GetUserId(c)
	if err != nil {
		return
	}

	sent, recieved, inventory, coins, err := h.service.Information.GetInformation(id)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, responses.InfoResponse{
		Coins:     coins,
		Inventory: inventory,
		CoinHistory: responses.CoinHistory{
			Recieved: recieved,
			Sent:     sent,
		},
	})
}
