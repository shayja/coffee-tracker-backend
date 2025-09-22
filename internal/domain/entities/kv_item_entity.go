// file: internal/domain/entities/kv_item_entity.go
package entities

type KVItem struct {
    Key   int `json:"key"`
    Value string  `json:"value"`
}
