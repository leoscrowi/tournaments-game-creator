package models

import (
	"github.com/google/uuid"
	"time"
)

type Game struct {
	GameID     uuid.UUID `json:"game_id"`
	GameStart  time.Time `json:"game_start"`
	GameTypeID uuid.UUID `json:"game_type_id"`
}

type GameType struct {
	GameTypeID   uuid.UUID `json:"game_type_id"`
	PlatformName string    `json:"platform_name"`
}
