package financing

import "errors"

var (
	ErrNotFound      = errors.New("aggregate not found")
	ErrAlreasyExists = errors.New("aggregate already exists")
)
