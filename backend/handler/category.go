package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/P-SEN371-Group-3/educartion/model"
	"github.com/P-SEN371-Group-3/educartion/rpc"
	categoryService "github.com/P-SEN371-Group-3/educartion/service/category"
)

func HandleCategories(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetAllCategories(w, r)
	case http.MethodPost:
		handleInsertCategory(w, r)
	default:
		rpc.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func HandleCategory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetCategoryById(w, r)
	case http.MethodPatch:
		handleUpdateCategory(w, r)
	case http.MethodDelete:
		handleDeleteCategory(w, r)
	default:
		rpc.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func HandleCategoryProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		rpc.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	handleGetProductsByCategory(w, r)
}

func getCategoryID(r *http.Request) (int, error) {
	id := r.URL.Query().Get("id")
	if id == "" {
		id = r.PathValue("id")
	}

	if id == "" {
		path := r.URL.Path
		switch {
		case strings.HasPrefix(path, "/api/categories/update/"):
			id = strings.TrimPrefix(path, "/api/categories/update/")
		case strings.HasPrefix(path, "/api/categories/delete/"):
			id = strings.TrimPrefix(path, "/api/categories/delete/")
		case strings.HasPrefix(path, "/api/categories/"):
			id = strings.TrimPrefix(path, "/api/categories/")
		}
	}

	return strconv.Atoi(id)
}

func handleGetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := categoryService.GetAllCategories()
	if err != nil {
		rpc.WriteError(w, http.StatusInternalServerError, "failed to retrieve categories")
		return
	}

	rpc.WriteJSON(w, http.StatusOK, categories)
}

func handleGetCategoryById(w http.ResponseWriter, r *http.Request) {
	categoryID, err := getCategoryID(r)
	if err != nil {
		rpc.WriteError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	category, err := categoryService.GetCategoryById(categoryID)
	if err == categoryService.ErrCategoryNotFound {
		rpc.WriteError(w, http.StatusNotFound, "category not found")
		return
	}
	if err != nil {
		rpc.WriteError(w, http.StatusInternalServerError, "failed to retrieve category")
		return
	}

	rpc.WriteJSON(w, http.StatusOK, category)
}

func handleGetProductsByCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := getCategoryID(r)
	if err != nil {
		rpc.WriteError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	products, err := categoryService.GetProductsByCategory(categoryID)
	if err != nil {
		rpc.WriteError(w, http.StatusInternalServerError, "failed to retrieve products for category")
		return
	}

	rpc.WriteJSON(w, http.StatusOK, products)
}

func handleInsertCategory(w http.ResponseWriter, r *http.Request) {
	var req model.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rpc.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		rpc.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}

	category, err := categoryService.CreateCategory(req)
	if err == categoryService.ErrCategoryNameRequired {
		rpc.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}
	if err != nil {
		rpc.WriteError(w, http.StatusInternalServerError, "failed to create category")
		return
	}

	rpc.WriteJSON(w, http.StatusCreated, category)
}

func handleUpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := getCategoryID(r)
	if err != nil {
		rpc.WriteError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	var req model.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		rpc.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		rpc.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}

	updatedCategory, err := categoryService.UpdateCategory(categoryID, req)
	if err == categoryService.ErrCategoryNotFound {
		rpc.WriteError(w, http.StatusNotFound, "category not found")
		return
	}
	if err == categoryService.ErrCategoryNameRequired {
		rpc.WriteError(w, http.StatusBadRequest, "name is required")
		return
	}
	if err != nil {
		rpc.WriteError(w, http.StatusInternalServerError, "failed to update category")
		return
	}

	rpc.WriteJSON(w, http.StatusOK, updatedCategory)
}

func handleDeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID, err := getCategoryID(r)
	if err != nil {
		rpc.WriteError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	err = categoryService.DeleteCategory(categoryID)
	if err == categoryService.ErrCategoryNotFound {
		rpc.WriteError(w, http.StatusNotFound, "category not found")
		return
	}
	if err != nil {
		rpc.WriteError(w, http.StatusInternalServerError, "failed to delete category")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
