package handler

import (
	"errors"
	"net/http"
	"strings"

	"avito_go/pkg/responses"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
)

func (h *Handler) UserIdentity(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		responses.NewErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		responses.NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}

	if headerParts[0] != "Bearer" {
		responses.NewErrorResponse(c, http.StatusUnauthorized, "invalid auth header, not Bearer token")
		return
	}

	if headerParts[1] == "" {
		responses.NewErrorResponse(c, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, coins, err := h.service.Authorization.ParseToken(headerParts[1])

	if err != nil {
		responses.NewErrorResponse(c, http.StatusUnauthorized, err.Error())
		return
	}
	c.Set("userId", userId)
	c.Set("coins", coins)
}

func GetUserId(c *gin.Context) (int, error) {
	id, ok := c.Get("userId")
	if !ok {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}

	idInt, ok := id.(int)
	if !ok {
		responses.NewErrorResponse(c, http.StatusInternalServerError, "cannot convert to int")
		return 0, errors.New("cannot convert to int")
	}

	return idInt, nil
}
