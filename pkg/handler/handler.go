package handler

import (
	"avito_go/pkg/service"

	"github.com/gin-gonic/gin"
)

var validItems = map[string]bool{
	"t-shirt":    true,
	"cup":        true,
	"book":       true,
	"pen":        true,
	"powerbank":  true,
	"hoody":      true,
	"umbrella":   true,
	"socks":      true,
	"wallet":     true,
	"pink-hoody": true,
}

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	routes1 := router.Group("/api", h.UserIdentity)
	{
		routes1.GET("/info", h.GetInformation)
		routes1.POST("/sendCoin", h.SendCoin)
		routes1.GET("/buy/:item", h.BuyItem)
	}

	routes2 := router.Group("/api")
	{
		routes2.POST("/auth", h.Auth)
	}
	return router
}
