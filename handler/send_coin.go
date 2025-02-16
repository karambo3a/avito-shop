package handler

import (
	"avito_go/requests"
	"avito_go/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) sendCoin(c *gin.Context) {
	id, err := GetUserId(c)
	if err != nil {
		return
	}

	var sendCoinRequest requests.SendCoinRequest

	if err := c.BindJSON(&sendCoinRequest); err != nil {
		responses.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	_, err = h.service.SendCoin.Send(id, sendCoinRequest.ToUser, sendCoinRequest.Amount)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusOK)
}
