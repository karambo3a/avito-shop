package handler

import (
	"avito_go/pkg/handler"
	"avito_go/pkg/requests"
	"avito_go/pkg/service"
	mock_service "avito_go/pkg/service/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_sendCoin(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSendCoinService := mock_service.NewMockSendCoin(ctrl)
	mockService := &service.Service{SendCoin: mockSendCoinService}

	handler := handler.NewHandler(mockService)

	tests := []struct {
		name           string
		input          requests.SendCoinRequest
		headerName     string
		headerValue    string
		userId         any
		mockBehavior   func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			input: requests.SendCoinRequest{
				ToUser: "testuser",
				Amount: 100,
			},
			headerName:  "Authorization",
			headerValue: "Bearer token",
			userId:      1,
			mockBehavior: func() {
				mockSendCoinService.EXPECT().
					Send(1, "testuser", 100).Return(1, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   ``,
		},
		{
			name: "Bad request",
			input: requests.SendCoinRequest{
				ToUser: "",
				Amount: 100,
			},
			headerName:  "Authorization",
			headerValue: "Bearer token",
			userId:      1,
			mockBehavior: func() {
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"Key: 'SendCoinRequest.ToUser' Error:Field validation for 'ToUser' failed on the 'required' tag"}`,
		},
		{
			name: "Service Error",
			input: requests.SendCoinRequest{
				ToUser: "recipient",
				Amount: 100,
			},
			headerName:  "Authorization",
			headerValue: "Bearer token",
			userId:      1,
			mockBehavior: func() {
				mockSendCoinService.EXPECT().
					Send(1, "recipient", 100).
					Return(0, errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"service error"}`,
		},
		{
			name: "User not found",
			input: requests.SendCoinRequest{
				ToUser: "recipient",
				Amount: 100,
			},
			headerName:  "Authorization",
			headerValue: "Bearer token",
			userId:      1,
			mockBehavior: func() {
				mockSendCoinService.EXPECT().
					Send(1, "recipient", 100).
					Return(0, errors.New("user id not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"user id not found"}`,
		},
		{
			name: "Not enough coins",
			input: requests.SendCoinRequest{
				ToUser: "recipient",
				Amount: 10000,
			},
			headerName:  "Authorization",
			headerValue: "Bearer token",
			userId:      1,
			mockBehavior: func() {
				mockSendCoinService.EXPECT().
					Send(1, "recipient", 10000).
					Return(0, errors.New("not enough coins"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"not enough coins"}`,
		},
		{
			name: "Get user id error",
			input: requests.SendCoinRequest{
				ToUser: "recipient",
				Amount: 10000,
			},
			headerName:     "Authorization",
			headerValue:    "Bearer token",
			userId:         "a",
			mockBehavior:   func() {},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"cannot convert to int"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior()

			router := gin.New()

			router.POST("/sendCoin", func(c *gin.Context) {
				c.Set("userId", test.userId)
			}, handler.SendCoin)

			body := `{"toUser":"` + test.input.ToUser + `","amount":` + strconv.Itoa(test.input.Amount) + `}`
			req, _ := http.NewRequest(http.MethodPost, "/sendCoin", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set(test.headerName, test.headerValue)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}
