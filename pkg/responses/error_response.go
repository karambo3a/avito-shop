package responses

import (
	"log"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Errors string `json:"errors"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	log.Default().Println(message)
	c.AbortWithStatusJSON(statusCode, ErrorResponse{message})
}
