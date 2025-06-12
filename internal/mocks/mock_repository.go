// File: mocks/mock_repository.go
package mocks

import (
	"github.com/m-garey/fetchit-backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateTables() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRepository) InsertUser(req models.UserRequest) (models.UserResponse, error) {
	args := m.Called(req)
	return args.Get(0).(models.UserResponse), args.Error(1)
}

func (m *MockRepository) InsertStore(req models.StoreRequest) (models.StoreResponse, error) {
	args := m.Called(req)
	return args.Get(0).(models.StoreResponse), args.Error(1)
}

func (m *MockRepository) UpsertStar(req models.PurchaseRequest) (models.PurchaseResponse, error) {
	args := m.Called(req)
	return args.Get(0).(models.PurchaseResponse), args.Error(1)
}

func (m *MockRepository) GetSticker(userID, storeID string) (models.UserStickerResponse, error) {
	args := m.Called(userID, storeID)
	return args.Get(0).(models.UserStickerResponse), args.Error(1)
}

func (m *MockRepository) GetStickersByUser(userID string) (models.StickerByUserResponse, error) {
	args := m.Called(userID)
	return args.Get(0).(models.StickerByUserResponse), args.Error(1)
}
