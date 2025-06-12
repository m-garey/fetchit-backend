package handler

import (
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

// @Summary Create a new user
// @Description Create a user with a given username
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.UserRequest true "User info"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var req models.UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.repository.InsertUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert user"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Create a new store
// @Description Register a new store with name and location
// @Tags Stores
// @Accept json
// @Produce json
// @Param store body models.StoreRequest true "Store info"
// @Success 200 {object} models.StoreResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/stores [post]
func (h *Handler) CreateStore(c *gin.Context) {
	var req models.StoreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.repository.InsertStore(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to insert store"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Record a user purchase
// @Description Record a purchase and potentially award or level up a sticker
// @Tags Purchases
// @Accept json
// @Produce json
// @Param purchase body models.PurchaseRequest true "Purchase info"
// @Success 200 {object} models.PurchaseResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/purchase [post]
func (h *Handler) RecordPurchase(c *gin.Context) {
	var req models.PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.repository.UpsertStar(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update sticker progress"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Get a specific user-store sticker
// @Description Retrieve a sticker for a given user and store
// @Tags Stickers
// @Produce json
// @Param user_id query string true "User ID"
// @Param store_id query string true "Store ID"
// @Success 200 {object} models.UserStickerResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/stickers/:user_id/:store_id [get]
func (h *Handler) GetSticker(c *gin.Context) {
	userID := c.Query("user_id")
	storeID := c.Query("store_id")

	resp, err := h.repository.GetSticker(userID, storeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get sticker"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

// @Summary Get all stickers for a user
// @Description Retrieve all stickers that belong to a specific user
// @Tags Stickers
// @Produce json
// @Param user_id query string true "User ID"
// @Success 200 {object} models.StickerByUserResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/stickers/:user_id [get]
func (h *Handler) GetStickersByUser(c *gin.Context) {
	userID := c.Query("user_id")

	resp, err := h.repository.GetStickersByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user stickers"})
		return
	}

	c.JSON(http.StatusOK, resp)
}
