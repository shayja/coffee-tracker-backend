// file: internal/infrastructure/http/models/kv_dto.go
package models

import "coffee-tracker-backend/internal/entities"

type KVResponse struct {
	Items []entities.KVItem `json:"items"`
}