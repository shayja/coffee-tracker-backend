// file: internal/infrastructure/repositories/generic_kv_repository_impl.go
package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"coffee-tracker-backend/internal/entities"
	"coffee-tracker-backend/internal/repositories"

	"github.com/patrickmn/go-cache"
)

type GenericKVRepositoryImpl struct {
    db    *sql.DB
    cache *cache.Cache
}

func NewGenericKVRepositoryImpl(db *sql.DB) repositories.GenericKVRepository {
    c := cache.New(10*time.Minute, 15*time.Minute) // 10 min TTL, 15 min cleanup
    return &GenericKVRepositoryImpl{db: db, cache: c}
}

func (r *GenericKVRepositoryImpl) GetKV(ctx context.Context, typeID int, languageCode string) ([]entities.KVItem, error) {
    cacheKey := fmt.Sprintf("type:%d:lang:%s", typeID, languageCode)

    // Try cache first
    if cached, found := r.cache.Get(cacheKey); found {
        if items, ok := cached.([]entities.KVItem); ok {
            return items, nil
        }
        // If type assertion fails fallback to DB query below
    }

    var query string
    switch typeID {
    case 1: // coffee_types
        query = `
            SELECT ct.id, ctt.name
            FROM coffee_types ct
            JOIN coffee_type_translations ctt ON ct.id = ctt.coffee_type_id
            JOIN languages l ON ctt.language_id = l.id
            WHERE l.code = $1
            ORDER BY ct.order_by ASC
        `
    case 2: // sizes
        query = `
            SELECT s.id, st.name
            FROM coffee_sizes s
            JOIN coffee_size_translations st ON s.id = st.coffee_size_id
            JOIN languages l ON st.language_id = l.id
            WHERE l.code = $1
            ORDER BY s.order_by ASC
        `
    default:
        return nil, fmt.Errorf("unsupported typeID: %d", typeID)
    }

    rows, err := r.db.QueryContext(ctx, query, languageCode)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var items []entities.KVItem
    for rows.Next() {
        var item entities.KVItem
        if err := rows.Scan(&item.Key, &item.Value); err != nil {
            return nil, err
        }
        items = append(items, item)
    }

    // Cache the results before returning
    r.cache.Set(cacheKey, items, cache.DefaultExpiration)
    return items, nil
}
