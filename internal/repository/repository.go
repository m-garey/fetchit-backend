package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/m-garey/fetchit-backend/internal/models"
)

type Repository struct {
	conn *pgx.Conn
}

type API interface {
	CreateTables() error
	InsertUser(models.UserRequest) (models.UserResponse, error)
	InsertStore(models.StoreRequest) (models.StoreResponse, error)
	UpsertStar(models.PurchaseRequest) (models.PurchaseResponse, error)
	GetSticker(string, string) (models.UserStickerResponse, error)
	GetStickersByUser(string) (models.StickerByUserResponse, error)
}

func New(db *pgx.Conn) *Repository {
	return &Repository{conn: db}
}

func (r *Repository) CreateTables() error {
	schema := `
	CREATE TABLE Users (
	user_id UUID PRIMARY KEY,
	username VARCHAR(50) NOT NULL,
	email VARCHAR(100),
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE Stores (
	store_id UUID PRIMARY KEY,
	store_name VARCHAR(100) NOT NULL,
	location VARCHAR(255),
	sticker_theme VARCHAR(100),
	is_active BOOLEAN DEFAULT TRUE
	);

	CREATE TABLE User_Sticker_Progress (
	user_sticker_id UUID PRIMARY KEY,
	user_id UUID REFERENCES Users(user_id),
	store_id UUID REFERENCES Stores(store_id),
	current_level VARCHAR(20) CHECK (current_level IN ('bronze', 'silver', 'gold', 'platinum')),
	star_count INT DEFAULT 0,
	last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE Purchases (
	purchase_id UUID PRIMARY KEY,
	user_id UUID REFERENCES Users(user_id),
	store_id UUID REFERENCES Stores(store_id),
	purchase_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	source VARCHAR(50)
	);

	CREATE TABLE Sticker_Level_Requirements (
	level VARCHAR(20) PRIMARY KEY,
	stars_required INT,
	next_level VARCHAR(20)
	);
	`
	_, err := r.conn.Exec(context.Background(), schema)
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) InsertUser(user models.UserRequest) (models.UserResponse, error) {
	var id string
	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO users (username) VALUES ($1) ON CONFLICT (username) DO UPDATE SET username=EXCLUDED.username RETURNING id`, user.Username).Scan(&id)
	if err != nil {
		return models.UserResponse{}, err
	}
	return models.UserResponse{
		ID: id,
	}, nil
}

func (r *Repository) InsertStore(store models.StoreRequest) (models.StoreResponse, error) {
	var id string
	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO stores (name, location) VALUES ($1, $2) RETURNING id`, store.Name, store.Location).Scan(&id)
	if err != nil {
		return models.StoreResponse{}, err
	}
	return models.StoreResponse{
		ID: id,
	}, nil
}

func (r *Repository) UpsertStar(purchase models.PurchaseRequest) (models.PurchaseResponse, error) {
	var stars int
	var level string
	var level_up bool

	// Insert sticker if not exists
	_, _ = r.conn.Exec(context.Background(),
		`INSERT INTO User_Sticker_Progress (user_id, store_id) VALUES ($1, $2)
		ON CONFLICT DO NOTHING`, purchase.UserID, purchase.StoreID)

	err := r.conn.QueryRow(context.Background(),
		`SELECT stars, level FROM stickers WHERE user_id=$1 AND store_id=$2`, purchase.UserID, purchase.StoreID).Scan(&stars, &level)
	if err != nil {
		return models.PurchaseResponse{}, err
	}

	stars++
	newLevel := level
	switch {
	case level == "bronze" && stars >= 5:
		level_up = true
		newLevel = "silver"
	case level == "silver" && stars >= 15:
		level_up = true
		newLevel = "gold"
	default:
		level_up = false
	}

	_, err = r.conn.Exec(context.Background(),
		`UPDATE User_Sticker_Progress SET stars=$1, level=$2 WHERE user_id=$3 AND store_id=$4`, stars, newLevel, purchase.UserID, purchase.StoreID)
	if err != nil {
		return models.PurchaseResponse{}, err
	}

	return models.PurchaseResponse{
		LevelUp:   level_up,
		Level:     newLevel,
		StarCount: stars,
	}, nil
}

func (r *Repository) GetSticker(userID string, storeID string) (models.UserStickerResponse, error) {
	var stars int
	var level string
	var location string
	var storeName string

	err := r.conn.QueryRow(context.Background(),
		`SELECT st.store_name, st.location, usp.star_count, usp.current_level
		 FROM User_Sticker_Progress usp
		 JOIN Stores st ON usp.store_id = st.store_id
		 WHERE usp.user_id = $1 AND usp.store_id = $2`, userID, storeID).Scan(&storeName, &location, &stars, &level)
	if err != nil {
		return models.UserStickerResponse{}, err
	}

	return models.UserStickerResponse{
		StoreName: storeName,
		Location:  location,
		StarCount: stars,
		Level:     level,
	}, nil
}

func (r *Repository) GetStickersByUser(userID string) (models.StickerByUserResponse, error) {
	var resp models.StickerByUserResponse

	rows, err := r.conn.Query(context.Background(),
		`SELECT s.id, st.name, st.location, s.stars, s.level FROM User_Sticker_Progress s
		JOIN stores st ON s.store_id = st.id WHERE s.user_id = $1`, userID)
	if err != nil {
		return models.StickerByUserResponse{}, err
	}

	for rows.Next() {
		var id int
		var storeName, location, level string
		var stars int

		err := rows.Scan(&id, &storeName, &location, &stars, &level)
		if err != nil {
			return models.StickerByUserResponse{}, err
		}
		sticker := models.UserStickerResponse{
			StoreName: storeName,
			Location:  location,
			StarCount: stars,
			Level:     level,
		}

		resp.Stickers = append(resp.Stickers, sticker)

		fmt.Printf("Sticker %d - %s (%s): %d stars, Level: %s\n", id, storeName, location, stars, level)
	}

	return resp, nil
}
