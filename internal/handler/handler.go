package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/m-garey/fetchit-backend/internal/models"
	"github.com/m-garey/fetchit-backend/internal/repository"
)

type Handler struct {
	repository repository.API
}

type API interface {
	CreateUser(c *gin.Context)
	CreateStore(c *gin.Context)
	RecordPurchase(c *gin.Context)
	GetSticker(c *gin.Context)
	GetStickersByUser(c *gin.Context)
}

func New(repository repository.API) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) CreateUser(c *gin.Context) {
	var req models.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	id, err := h.repository.InsertUser(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert user"})
		return
	}

	c.JSON(http.StatusOK, models.UserResponse{ID: fmt.Sprintf("%d", id)})
}

func (h *Handler) CreateStore(c *gin.Context) {
	var req models.StoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	id, err := h.repository.InsertStore(req.Name, req.Location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert store"})
		return
	}

	c.JSON(http.StatusOK, models.StoreResponse{ID: fmt.Sprintf("%d", id)})
}

func (h *Handler) RecordPurchase(c *gin.Context) {
	var req models.PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.repository.UpsertStar(toInt(req.UserID), toInt(req.StoreID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update sticker progress"})
		return
	}

	c.JSON(http.StatusOK, models.PurchaseResponse{LevelUp: true, Level: "silver", StarCount: 5})
}

func (h *Handler) GetSticker(c *gin.Context) {
	userID := c.Query("user_id")
	storeID := c.Query("store_id")

	err := h.repository.GetSticker(toInt(userID), toInt(storeID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sticker"})
		return
	}

	c.Status(http.StatusOK)
}

func (h *Handler) GetStickersByUser(c *gin.Context) {
	userID := c.Query("user_id")

	err := h.repository.GetStickersByUser(toInt(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user stickers"})
		return
	}

	c.Status(http.StatusOK)
}

func toInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}
