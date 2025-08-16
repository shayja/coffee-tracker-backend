package repositories

import "errors"

var (
    // ErrNotFound is returned when a record is not found in the repository.
    ErrNotFound = errors.New("item not found")
)
