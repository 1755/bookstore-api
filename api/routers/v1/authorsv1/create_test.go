package authorsv1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/1755/bookstore-api/api/routers"
	"github.com/1755/bookstore-api/api/routers/v1/authorsv1"
	"github.com/1755/bookstore-api/internal/author"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestCreateTestSuite(t *testing.T) {
	suite.Run(t, new(AddressTestSuccessSuite))
}

type AddressTestSuccessSuite struct {
	suite.Suite
	ctx           context.Context
	config        *routers.Config
	authorService *author.MockService
	router        *gin.Engine

	builder *authorsv1.CreateRouterBuilder
}

func (s *AddressTestSuccessSuite) SetupTest() {
	gin.SetMode(gin.TestMode)

	s.ctx = context.Background()
	s.config = &routers.Config{
		BaseURL: "http://localhost:8080",
	}
	s.authorService = author.NewMockService(s.T())

	s.builder = authorsv1.NewCreateRouterBuilder(s.authorService, s.config)
	s.router = gin.Default()
	s.builder.Build(s.router.Group("/create"))
}

func (s *AddressTestSuccessSuite) TestCreateSuccess() {
	// Arrange
	now := time.Now()

	s.authorService.EXPECT().Create(
		s.ctx,
		&author.Model{
			Name: "The Ancient Dude",
			Bio:  "Very old author of the ancient times",
		},
	).Return(&author.Model{
		ID:        author.ID(100500),
		Name:      "The Ancient Dude",
		Bio:       "Very old author of the ancient times",
		CreatedAt: now,
		UpdatedAt: now,
	}, nil)

	requestBody, err := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"attributes": map[string]interface{}{
				"bio":  "Very old author of the ancient times",
				"name": "The Ancient Dude",
			},
			"type": "authors",
		},
	})
	if err != nil {
		s.T().Fatalf("Failed to marshal request data: %v", err)
	}

	expectedBody, err := json.Marshal(map[string]interface{}{
		"links": map[string]interface{}{
			"self": fmt.Sprintf("%s/v1/authors/100500", s.config.BaseURL),
		},
		"data": map[string]interface{}{
			"id":   "100500",
			"type": "authors",
			"attributes": map[string]interface{}{
				"bio":       "Very old author of the ancient times",
				"name":      "The Ancient Dude",
				"createdAt": now.Format(time.RFC3339Nano),
				"updatedAt": now.Format(time.RFC3339Nano),
			},
		},
	})
	if err != nil {
		s.T().Fatalf("Failed to marshal request data: %v", err)
	}

	// Act
	req, err := http.NewRequest(http.MethodPost, "/create/", bytes.NewBuffer(requestBody))
	if err != nil {
		s.T().Fatalf("Failed to create request: %v", err)
	}
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	// Assett
	assert.Equal(s.T(), http.StatusCreated, w.Code)
	assert.JSONEq(s.T(), string(expectedBody), w.Body.String())
}
