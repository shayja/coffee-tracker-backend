package dto

import "coffee-tracker-backend/internal/domain/entities"

type KVResponse struct {
	Items []entities.KVItem `json:"items"`
}