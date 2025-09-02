package models

import "github.com/google/uuid"

type Result struct {
	ResultID uuid.UUID `json:"result_id"`
	GameID   uuid.UUID `json:"game_id"`
	WinnerID uuid.UUID `json:"winner_id"`
	Comment  string    `json:"comment"`
}
