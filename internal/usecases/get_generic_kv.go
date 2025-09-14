// file: internal/usecases/get_generic_kv.go
package usecases

import (
	"coffee-tracker-backend/internal/domain/entities"
	"coffee-tracker-backend/internal/domain/repositories"
	"context"
)

type GetGenericKVUseCase struct {
    genericKvRepo repositories.GenericKVRepository
}

func NewGetGenericKVUseCase(genericKvRepo repositories.GenericKVRepository) *GetGenericKVUseCase {
	return &GetGenericKVUseCase{genericKvRepo: genericKvRepo}
}

func (uc *GetGenericKVUseCase) Execute(ctx context.Context, typeID int, languageCode string) ([]entities.KVItem, error) {
    return uc.genericKvRepo.GetKV(ctx, typeID, languageCode)
}
