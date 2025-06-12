package models

import "time"

// USER

type User struct {
	ID        string    `json:"user_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRequest struct {
	Username string `json:"username"`
}

type UserResponse struct {
	ID string `json:"user_id"`
}

// STORE

type Store struct {
	ID           string `json:"store_id"`
	Name         string `json:"name"`
	Location     string `json:"location"`
	StickerTheme string `json:"sticker_theme"`
	IsActive     bool   `json:"is_active"`
}

type StoreRequest struct {
	Name     string `json:"store_name"`
	Location string `json:"location"`
}

type StoreResponse struct {
	ID string `json:"store_id"`
}

// STICKER

type UserStickerProgress struct {
	ID           string    `json:"sticker_id"`
	UserID       string    `json:"user_id"`
	StoreID      string    `json:"store_id"`
	CurrentLevel string    `json:"current_level"`
	StarCount    int       `json:"star_count"`
	LastUpdated  time.Time `json:"last_updated"`
}

type Purchase struct {
	ID           string    `json:"purchase_id"`
	UserID       string    `json:"user_id"`
	StoreID      string    `json:"store_id"`
	PurchaseTime time.Time `json:"purchase_time"`
}

type PurchaseRequest struct {
	UserID  string `json:"user_id"`
	StoreID string `json:"store_id"`
}

type PurchaseResponse struct {
	LevelUp   bool   `json:"level_up"`
	Level     string `json:"level"`
	StarCount int    `json:"star_count"`
}

type StickerLevelRequirement struct {
	Level         string `json:"level"`
	StarsRequired int    `json:"stars_required"`
	NextLevel     string `json:"next_level"`
}

// Get sticker for user for specific store

type UserStickerResponse struct {
	StoreName string `json:"store_name"`
	Location  string `json:"location"`
	StarCount int    `json:"star_count"`
	Level     string `json:"level"`
}

type StickerByUserResponse struct {
	Stickers []UserStickerResponse `json:"stickers"`
}
