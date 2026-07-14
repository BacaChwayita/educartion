package category

import "errors"

var (
	ErrCategoryNotFound     = errors.New("category not found")
	ErrCategoryNameRequired = errors.New("category name is required")
)
