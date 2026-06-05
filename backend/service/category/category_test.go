package category

import (
	"testing"

	"github.com/P-SEN371-Group-3/educartion/model"
)

func TestCreateCategory_EmptyNameReturnsErrCategoryNameRequired(t *testing.T) {
	_, err := CreateCategory(model.CreateCategoryRequest{
		Name:        "",
		Description: "A category without a name",
		Is_active:   true,
	})

	if err != ErrCategoryNameRequired {
		t.Fatalf("expected ErrCategoryNameRequired, got %v", err)
	}
}

func TestUpdateCategory_EmptyNameReturnsErrCategoryNameRequired(t *testing.T) {
	_, err := UpdateCategory(1, model.UpdateCategoryRequest{
		Name:        "",
		Description: "A category without a name",
		Is_active:   true,
	})

	if err != ErrCategoryNameRequired {
		t.Fatalf("expected ErrCategoryNameRequired, got %v", err)
	}
}
