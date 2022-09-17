package get_test

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jeffersonto/feira-api/dto"
	"github.com/jeffersonto/feira-api/entity"
	"github.com/jeffersonto/feira-api/handlers"
	"github.com/jeffersonto/feira-api/handlers/get"
	"github.com/jeffersonto/feira-api/server/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFairByQuery(t *testing.T) {

	var service *serviceMock

	tests := []struct {
		name           string
		queryParameter string
		warmUP         func()
		expected       func(result *httptest.ResponseRecorder)
	}{
		{
			name:           "Should successfully get by query and return status code 200",
			queryParameter: "?distrito=VILA FORMOSA",
			warmUP: func() {
				service = new(serviceMock)
				service.On("FindFairByQuery", mock.Anything).Return([]entity.Fair{{ID: 1, NomeFeira: "Feira Teste"}}, nil)
			},
			expected: func(result *httptest.ResponseRecorder) {
				assert.Equal(t, 200, result.Code)
				service.AssertNumberOfCalls(t, "FindFairByQuery", 1)
			},
		},
		{
			name:           "Should execute the FindFairByQuery Function, however receive an internal_server_error with status code 500",
			queryParameter: "?distrito=VILA FORMOSA",
			warmUP: func() {
				service = new(serviceMock)
				service.On("FindFairByQuery", mock.Anything).Return(entity.Fair{}, fmt.Errorf("internal_server_error"))
			},
			expected: func(result *httptest.ResponseRecorder) {
				assert.Equal(t, 500, result.Code)
				service.AssertNumberOfCalls(t, "FindFairByQuery", 1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.warmUP()
			router := gin.Default()
			router.Use(middleware.ErrorHandle())
			handler := handlers.NewHandler(service)
			get.NewFairByQueryHandler(handler, router)
			response := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/feiras%v", tt.queryParameter), nil)
			router.ServeHTTP(response, req)
			tt.expected(response)
		})
	}
}

func (sm *serviceMock) FindFairByQuery(filters dto.QueryParameters) ([]entity.Fair, error) {
	args := sm.Called(filters)
	result, _ := args.Get(0).([]entity.Fair)
	return result, args.Error(1)
}
