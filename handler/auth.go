package handler

import (
	"avito_go/requests"
	"avito_go/responses"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) auth(c *gin.Context) {
	var input requests.AuthRequest

	if err := c.BindJSON(&input); err != nil {
		responses.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		responses.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, responses.AuthResponse{
		Token: token,
	})
}
