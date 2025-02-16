package handler

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"avito_go/responses"
	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
)

func (h *Handler) userIdentity(c *gin.Context) {
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

	userId, coins, err := h.service.Authorization.ParseToken(headerParts[1])
	log.Default().Println(userId)
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

	log.Default().Println(idInt)
	return idInt, nil
}
