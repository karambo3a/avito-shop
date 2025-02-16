package handler

import (
	"avito_go/pkg/handler"
	"avito_go/pkg/service"
	mock_service "avito_go/pkg/service/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_buyItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockBuyItemService := mock_service.NewMockBuyItem(ctrl)
	mockService := &service.Service{BuyItem: mockBuyItemService}

	handler := handler.NewHandler(mockService)

	tests := []struct {
		name           string
		headerName     string
		headerValue    string
		userId         any
		item           string
		mockBehavior   func(item string)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Success",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			userId:      1,
			item:        "cup",
			mockBehavior: func(item string) {
				mockBuyItemService.EXPECT().
					Buy(1, item).Return(1, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   ``,
		},
		{
			name:           "Invalid item",
			headerName:     "Authorization",
			headerValue:    "Bearer token",
			userId:         1,
			item:           "banana",
			mockBehavior:   func(item string) {},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"invalid item"}`,
		},
		{
			name:        "Not enough coins",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			userId:      1,
			item:        "cup",
			mockBehavior: func(item string) {
				mockBuyItemService.EXPECT().
					Buy(1, item).Return(0, errors.New("not enough coins"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"not enough coins"}`,
		},
		{
			name:           "Get user id error",
			headerName:     "Authorization",
			headerValue:    "Bearer token",
			userId:         "a",
			item:           "cup",
			mockBehavior:   func(item string) {},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"cannot convert to int"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.item)

			router := gin.New()

			router.GET("/buy/:item", func(c *gin.Context) {
				c.Set("userId", test.userId)
			}, handler.BuyItem)

			body := ""
			req, _ := http.NewRequest(http.MethodGet, "/buy/"+test.item, strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set(test.headerName, test.headerValue)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}
