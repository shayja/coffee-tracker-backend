package repositories

import (
	"coffee-tracker-backend/internal/domain/entities"
	"context"
)


type GenericKVRepository interface {
	GetKV(ctx context.Context, typeID int, languageCode string) ([]entities.KVItem, error)
}
