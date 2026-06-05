package category

import (
	"errors"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/service/db"
	"github.com/jackc/pgx/v5"
)

func GetAllCategories() ([]model.Category, error) {
	return db.GetAllCategories()
}

func GetCategoryById(categoryID int) (model.Category, error) {
	category, err := db.GetCategoryById(categoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Category{}, ErrCategoryNotFound
		}
		return model.Category{}, err
	}

	return category, nil
}

func CreateCategory(req model.CreateCategoryRequest) (model.Category, error) {
	if req.Name == "" {
		return model.Category{}, ErrCategoryNameRequired
	}

	category := model.Category{
		Name:        req.Name,
		Description: req.Description,
		Is_active:   req.Is_active,
	}

	return db.InsertCategory(category)
}

func UpdateCategory(categoryID int, req model.UpdateCategoryRequest) (model.Category, error) {
	if req.Name == "" {
		return model.Category{}, ErrCategoryNameRequired
	}

	category := model.Category{
		Category_id: categoryID,
		Name:        req.Name,
		Description: req.Description,
		Is_active:   req.Is_active,
	}

	rows, err := db.UpdateCategory(category)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Category{}, ErrCategoryNotFound
		}
		return model.Category{}, err
	}

	if rows == 0 {
		return model.Category{}, ErrCategoryNotFound
	}

	return db.GetCategoryById(categoryID)
}

func DeleteCategory(categoryID int) error {
	rows, err := db.DeleteCategoryById(categoryID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCategoryNotFound
		}
		return err
	}

	if rows == 0 {
		return ErrCategoryNotFound
	}

	return nil
}

func GetProductsByCategory(categoryID int) ([]model.Product, error) {
	return db.GetProductsByCategoryId(categoryID)
}
