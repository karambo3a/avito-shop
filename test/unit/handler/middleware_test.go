package handler

import (
	"avito_go/pkg/handler"
	"avito_go/pkg/service"
	mock_service "avito_go/pkg/service/mocks"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_userIdentity(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAuthService := mock_service.NewMockAuthorization(ctrl)

	mockService := &service.Service{Authorization: mockAuthService}
	handler := handler.NewHandler(mockService)

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         func(token string)
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "Ok",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(token string) {
				mockAuthService.EXPECT().ParseToken(token).Return(1, 1000, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "1",
		},
		{
			name:                 "Invalid Header Name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"errors":"empty auth header"}`,
		},
		{
			name:                 "Invalid Header Value",
			headerName:           "Authorization",
			headerValue:          "Bearr token",
			token:                "token",
			mockBehavior:         func(token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"errors":"invalid auth header, not Bearer token"}`,
		},
		{
			name:                 "Empty Token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"errors":"token is empty"}`,
		},
		{
			name:        "Parse Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(token string) {
				mockAuthService.EXPECT().ParseToken(token).Return(0, 0, errors.New("invalid token"))
			},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"errors":"invalid token"}`,
		},
		{
			name:                 "No token",
			headerName:           "Authorization",
			headerValue:          "Bearer",
			token:                "token",
			mockBehavior:         func(token string) {},
			expectedStatusCode:   http.StatusUnauthorized,
			expectedResponseBody: `{"errors":"invalid auth header"}`,
		},
	}

	for _, test := range testTable {
		t.Run(test.name, func(t *testing.T) {
			test.mockBehavior(test.token)

			router := gin.New()
			router.GET("/identity", handler.UserIdentity, func(c *gin.Context) {
				id, _ := c.Get("userId")
				c.String(200, "%d", id)
			})

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/identity", nil)
			req.Header.Set(test.headerName, test.headerValue)

			router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, w.Body.String(), test.expectedResponseBody)
		})
	}
}

func TestHandler_GetUserId(t *testing.T) {
	setCntx := func(id any) *gin.Context {
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		c.Set("userId", id)
		return c
	}

	tests := []struct {
		name    string
		context *gin.Context
		id      int
		err     error
	}{
		{
			name:    "Valid user ID",
			context: setCntx(1),
			id:      1,
			err:     nil,
		},
		{
			name: "User ID not found",
			context: func(id any) *gin.Context {
				c, _ := gin.CreateTestContext(httptest.NewRecorder())
				return c
			}(1),
			id:  0,
			err: errors.New("user id not found"),
		},
		{
			name:    "Invalid user ID type",
			context: setCntx("a"),
			id:      0,
			err:     errors.New("cannot convert to int"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var c *gin.Context
			if test.context != nil {
				c = test.context
			}

			id, err := handler.GetUserId(c)

			if test.err != nil {
				assert.Error(t, err)
				assert.Equal(t, test.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.id, id)
			}
		})
	}
}
