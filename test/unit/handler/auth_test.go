package handler

import (
	"avito_go/pkg/handler"
	"avito_go/pkg/requests"
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

func TestHandler_auth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)
	mockService := &service.Service{Authorization: mockAuthService}

	handler := handler.NewHandler(mockService)

	tests := []struct {
		name           string
		input          requests.AuthRequest
		mockBehavior   func()
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Success",
			input: requests.AuthRequest{
				Username: "testuser",
				Password: "testpassword",
			},
			mockBehavior: func() {
				mockAuthService.EXPECT().
					GenerateToken("testuser", "testpassword").
					Return("valid-token", nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"token":"valid-token"}`,
		},
		{
			name: "Service Error",
			input: requests.AuthRequest{
				Username: "testuser",
				Password: "testpassword",
			},
			mockBehavior: func() {
				mockAuthService.EXPECT().
					GenerateToken("testuser", "testpassword").
					Return("", errors.New("service error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"service error"}`,
		},
		{
			name: "Bad request",
			input: requests.AuthRequest{
				Username: "",
				Password: "testpassword",
			},
			mockBehavior: func() {

			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"errors":"Key: 'AuthRequest.Username' Error:Field validation for 'Username' failed on the 'required' tag"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior()

			router := gin.New()
			router.POST("/auth", handler.Auth)

			body := `{"username":"` + test.input.Username + `","password":"` + test.input.Password + `"}`
			req, _ := http.NewRequest(http.MethodPost, "/auth", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}
