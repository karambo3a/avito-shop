package handler

import (
	"avito_go/requests"
	"avito_go/service"
	mock_service "avito_go/service/mocks"
	"errors"
	"log"
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

	handler := NewHandler(mockService)

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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockBehavior()

			router := gin.New()
			router.POST("/auth", handler.auth)

			body := `{"username":"` + tt.input.Username + `","password":"` + tt.input.Password + `"}`
			req, _ := http.NewRequest(http.MethodPost, "/auth", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			log.Default().Println(w.Code, w.Body.String())
			assert.Equal(t, tt.expectedStatus, w.Code)
			assert.Equal(t, tt.expectedBody, w.Body.String())
		})
	}
}
