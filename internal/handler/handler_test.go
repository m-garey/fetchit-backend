// File: handler_test.go
package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/m-garey/fetchit-backend/internal/handler"
	"github.com/m-garey/fetchit-backend/internal/mocks"
	"github.com/m-garey/fetchit-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func performRequest(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	b, _ := json.Marshal(body)
	req, _ := http.NewRequest(method, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestCreateUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	h := handler.New(mockRepo)
	r := gin.Default()
	r.POST("/api/users", h.CreateUser)

	req := models.UserRequest{Username: "tester"}
	resp := models.UserResponse{ID: "abc123"}
	mockRepo.On("InsertUser", req).Return(resp, nil)

	w := performRequest(r, "POST", "/api/users", req)
	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestCreateStore(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	h := handler.New(mockRepo)
	r := gin.Default()
	r.POST("/api/stores", h.CreateStore)

	req := models.StoreRequest{Name: "Store A", Location: "123 Main"}
	resp := models.StoreResponse{ID: "store123"}
	mockRepo.On("InsertStore", req).Return(resp, nil)

	w := performRequest(r, "POST", "/api/stores", req)
	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestRecordPurchase(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	h := handler.New(mockRepo)
	r := gin.Default()
	r.POST("/api/purchase", h.RecordPurchase)

	req := models.PurchaseRequest{UserID: "user1", StoreID: "store1"}
	resp := models.PurchaseResponse{LevelUp: true, Level: "silver", StarCount: 5}
	mockRepo.On("UpsertStar", req).Return(resp, nil)

	w := performRequest(r, "POST", "/api/purchase", req)
	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetSticker(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	h := handler.New(mockRepo)
	r := gin.Default()
	r.GET("/api/stickers/:user_id/:store_id", h.GetSticker)

	userID := "user1"
	storeID := "store1"
	resp := models.UserStickerResponse{
		StoreName: "Store A",
		Location:  "123 Main",
		StarCount: 3,
		Level:     "bronze",
	}
	mockRepo.On("GetSticker", userID, storeID).Return(resp, nil)

	req := httptest.NewRequest("GET", "/api/stickers/user1/store1?user_id=user1&store_id=store1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestGetStickersByUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	h := handler.New(mockRepo)
	r := gin.Default()
	r.GET("/api/stickers/:user_id", h.GetStickersByUser)

	userID := "user1"
	resp := models.StickerByUserResponse{
		Stickers: []models.UserStickerResponse{
			{StoreName: "Store A", Location: "123 Main", StarCount: 5, Level: "silver"},
		},
	}
	mockRepo.On("GetStickersByUser", userID).Return(resp, nil)

	req := httptest.NewRequest("GET", "/api/stickers/user1?user_id=user1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockRepo.AssertExpectations(t)
}
