package handler

import (
	"avito_go/pkg/handler"
	"avito_go/pkg/responses"
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

func TestHandler_getInformation(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockInformatiomService := mock_service.NewMockInformation(ctrl)
	mockService := &service.Service{Information: mockInformatiomService}

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
				mockInformatiomService.EXPECT().
					GetInformation(1).Return(make([]responses.Transaction, 0), make([]responses.Transaction, 0), make([]responses.InventoryItem, 0), 1000, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"coins":1000,"inventory":[],"coinHistory":{"received":[],"sent":[]}}`,
		},
		{
			name:        "Server error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			userId:      1,
			item:        "cup",
			mockBehavior: func(item string) {
				mockInformatiomService.EXPECT().
					GetInformation(1).Return(nil, nil, nil, 0, errors.New("server error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"errors":"server error"}`,
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

			router.GET("/info", func(c *gin.Context) {
				c.Set("userId", test.userId)
			}, handler.GetInformation)

			body := ""
			req, _ := http.NewRequest(http.MethodGet, "/info", strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set(test.headerName, test.headerValue)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatus, w.Code)
			assert.Equal(t, test.expectedBody, w.Body.String())
		})
	}
}
