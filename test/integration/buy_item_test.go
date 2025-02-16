package test

import (
	"avito_go/pkg/handler"
	"avito_go/pkg/repository"
	"avito_go/pkg/service"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_buyItem(t *testing.T) {

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DATABASE_HOST"),
		Port:     os.Getenv("DATABASE_PORT"),
		Username: os.Getenv("DATABASE_USER"),
		Password: os.Getenv("DATABASE_PASSWORD"),
		DBName:   os.Getenv("DATABASE_NAME"),
		SSLMode:  "disable",
	})
	assert.NoError(t, err)
	defer db.Close()

	repository := repository.NewRepository(db)
	service := service.NewService(repository)
	handler := handler.NewHandler(service)

	body := map[string]interface{}{
		"username": "testuserBuy",
		"password": "testpasswordBuy",
	}
	jsonBody, _ := json.Marshal(body)

	router := gin.New()
	router.POST("/auth", handler.Auth)
	req, err := http.NewRequest(http.MethodPost, "/auth", bytes.NewBuffer(jsonBody))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	token, exists := response["token"]
	assert.True(t, exists)
	assert.NotEmpty(t, token)

	router.GET("/buy/:item", handler.UserIdentity, handler.BuyItem)

	requestBody := ""
	req, _ = http.NewRequest(http.MethodGet, "/buy/cup", strings.NewReader(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token.(string))

	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "", w.Body.String())
}
