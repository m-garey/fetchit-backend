package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type Repository struct {
	conn *pgx.Conn
}

type API interface {
	CreateTables() error
	InsertUser(string) (int, error)
	InsertStore(string, string) (int, error)
	UpsertStar(int, int) error
	GetSticker(int, int) error
	GetStickersByUser(int) error
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
		log.Fatal("Table creation failed:", err)
		return err
	}
	return nil
}

func (r *Repository) InsertUser(username string) (int, error) {
	var id int
	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO users (username) VALUES ($1) ON CONFLICT (username) DO UPDATE SET username=EXCLUDED.username RETURNING id`, username).Scan(&id)
	if err != nil {
		log.Fatal("Insert user failed:", err)
	}
	return id, nil
}

func (r *Repository) InsertStore(name string, location string) (int, error) {
	var id int
	err := r.conn.QueryRow(context.Background(),
		`INSERT INTO stores (name, location) VALUES ($1, $2) RETURNING id`, name, location).Scan(&id)
	if err != nil {
		log.Fatal("Insert store failed:", err)
	}
	return id, nil
}

func (r *Repository) UpsertStar(userID int, storeID int) error {
	var stars int
	var level string

	// Insert sticker if not exists
	_, _ = r.conn.Exec(context.Background(),
		`INSERT INTO stickers (user_id, store_id) VALUES ($1, $2)
		ON CONFLICT DO NOTHING`, userID, storeID)

	err := r.conn.QueryRow(context.Background(),
		`SELECT stars, level FROM stickers WHERE user_id=$1 AND store_id=$2`, userID, storeID).Scan(&stars, &level)
	if err != nil {
		log.Fatal("Query sticker failed:", err)
	}

	stars++
	newLevel := level
	switch {
	case level == "bronze" && stars >= 5:
		newLevel = "silver"
	case level == "silver" && stars >= 15:
		newLevel = "gold"
	}

	_, err = r.conn.Exec(context.Background(),
		`UPDATE stickers SET stars=$1, level=$2 WHERE user_id=$3 AND store_id=$4`, stars, newLevel, userID, storeID)
	if err != nil {
		log.Fatal("Update sticker failed:", err)
	}

	return nil
}

func (r *Repository) GetSticker(userID int, storeID int) error {
	var stars int
	var level string
	err := r.conn.QueryRow(context.Background(),
		`SELECT stars, level FROM stickers WHERE user_id=$1 AND store_id=$2`, userID, storeID).Scan(&stars, &level)
	if err != nil {
		log.Fatal("Fetching sticker info failed:", err)
	}
	//fmt.Sprintf("🏅 Sticker Info → Stars: %d, Level: %s\n", stars, level)
	return nil
}

func (r *Repository) GetStickersByUser(userID int) error {
	rows, err := r.conn.Query(context.Background(),
		`SELECT s.id, st.name, st.location, s.stars, s.level FROM stickers s
		JOIN stores st ON s.store_id = st.id WHERE s.user_id = $1`, userID)
	if err != nil {
		log.Fatal("Fetching all sticker info failed:", err)
	}

	for rows.Next() {
		var id int
		var storeName, location, level string
		var stars int

		err := rows.Scan(&id, &storeName, &location, &stars, &level)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Sticker %d - %s (%s): %d stars, Level: %s\n", id, storeName, location, stars, level)
	}

	return nil
}
