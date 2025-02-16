package handler

import (
	"avito_go/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) buyItem(c *gin.Context) {
	id, err := GetUserId(c)
	if err != nil {
		return
	}

	item := c.Param("item")

	if !validItems[item] {
		responses.NewErrorResponse(c, http.StatusBadRequest, "invalid item")
		return
	}

	_, err = h.service.BuyItem.Buy(id, item)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusOK)
}
