// file: internal/infrastructure/repositories/generic_kv_repository_impl.go
package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"coffee-tracker-backend/internal/entities"
	"coffee-tracker-backend/internal/repositories"

	"github.com/patrickmn/go-cache"
)

var (
	ErrUnsupportedType = errors.New("unsupported KV type")

	kvQueries = map[int]string{
            // coffee types
		1: `SELECT ct.id, ctt.name
			FROM coffee_types ct
			JOIN coffee_type_translations ctt ON ct.id = ctt.coffee_type_id
			JOIN languages l ON ctt.language_id = l.id
			WHERE l.code = $1
			ORDER BY ct.order_by ASC`,
            // coffee sizes
		2: `SELECT s.id, st.name
			FROM coffee_sizes s
			JOIN coffee_size_translations st ON s.id = st.coffee_size_id
			JOIN languages l ON st.language_id = l.id
			WHERE l.code = $1
			ORDER BY s.order_by ASC`,
	}
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

	// 1️⃣ Try cache first
	if cached, found := r.cache.Get(cacheKey); found {
		if items, ok := cached.([]entities.KVItem); ok {
			return items, nil
		}
	}

	// 2️⃣ Lookup query
	query, ok := kvQueries[typeID]
	if !ok {
		return nil, ErrUnsupportedType
	}

	// 3️⃣ Query DB
	rows, err := r.db.QueryContext(ctx, query, languageCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]entities.KVItem, 0)
	for rows.Next() {
		var item entities.KVItem
		if err := rows.Scan(&item.Key, &item.Value); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	// 4️⃣ Cache results
	r.cache.Set(cacheKey, items, cache.DefaultExpiration)

	return items, nil
}
