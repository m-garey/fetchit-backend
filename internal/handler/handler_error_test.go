package handler_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/m-garey/fetchit-backend/internal/handler"
	"github.com/m-garey/fetchit-backend/internal/mocks"
	"github.com/m-garey/fetchit-backend/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	h := handler.New(mockRepo)
	r := gin.Default()
	r.POST("/api/users", h.CreateUser)

	invalidJSON := []byte(`{"invalid"}`)
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateUser_DBFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	h := handler.New(mockRepo)
	r := gin.Default()
	r.POST("/api/users", h.CreateUser)

	reqBody := models.UserRequest{Username: "tester"}
	mockRepo.On("InsertUser", reqBody).Return(models.UserResponse{}, errors.New("DB fail"))

	w := performRequest(r, "POST", "/api/users", reqBody)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestCreateStore_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	h := handler.New(mockRepo)
	r := gin.Default()
	r.POST("/api/stores", h.CreateStore)

	invalidJSON := []byte(`{"wrong"}`)
	req, _ := http.NewRequest("POST", "/api/stores", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRecordPurchase_DBFailure(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	h := handler.New(mockRepo)
	r := gin.Default()
	r.POST("/api/purchase", h.RecordPurchase)

	reqBody := models.PurchaseRequest{UserID: "u1", StoreID: "s1"}
	mockRepo.On("UpsertStar", reqBody).Return(models.PurchaseResponse{}, errors.New("fail"))

	w := performRequest(r, "POST", "/api/purchase", reqBody)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetSticker_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	h := handler.New(mockRepo)
	r := gin.Default()
	r.GET("/api/stickers/:user_id/:store_id", h.GetSticker)

	mockRepo.On("GetSticker", "u1", "s1").Return(models.UserStickerResponse{}, errors.New("fetch error"))

	req := httptest.NewRequest("GET", "/api/stickers/u1/s1?user_id=u1&store_id=s1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestGetStickersByUser_DBError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockRepo := new(mocks.MockRepository)
	h := handler.New(mockRepo)
	r := gin.Default()
	r.GET("/api/stickers/:user_id", h.GetStickersByUser)

	mockRepo.On("GetStickersByUser", "u1").Return(models.StickerByUserResponse{}, errors.New("fetch fail"))

	req := httptest.NewRequest("GET", "/api/stickers/u1?user_id=u1", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
