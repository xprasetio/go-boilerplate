package category

import (
	"net/http"
	"strconv"

	categoryModel "boilerplate/internal/category/model"
	userModel "boilerplate/internal/user/model"
	"boilerplate/pkg/response"

	"github.com/labstack/echo/v4"
)

type CategoryHandler struct {
	categoryService *CategoryService
}

func NewCategoryHandler(categoryService *CategoryService) *CategoryHandler {
	return &CategoryHandler{
		categoryService: categoryService,
	}
}

func (h *CategoryHandler) Create(c echo.Context) error {
	user, exists := c.Get("user").(*userModel.User)
	if !exists {
		return response.Unauthorized(c, "unauthorized", nil)
	}

	var input categoryModel.CreateCategoryInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "invalid request payload", err)
	}

	category, err := h.categoryService.Create(input, user)
	if err != nil {
		return response.BadRequest(c, "failed to create category", err)
	}

	return response.Success(c, http.StatusCreated, "Category created successfully", category)
}

func (h *CategoryHandler) GetAll(c echo.Context) error {
	categories, err := h.categoryService.GetAll()
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, "failed to get categories", err)
	}

	return response.Success(c, http.StatusOK, "Categories retrieved successfully", categories)
}

func (h *CategoryHandler) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid category id", err)
	}

	category, err := h.categoryService.GetByID(uint(id))
	if err != nil {
		return response.BadRequest(c, "category not found", err)
	}

	return response.Success(c, http.StatusOK, "Category retrieved successfully", category)
}

func (h *CategoryHandler) Update(c echo.Context) error {
	user, exists := c.Get("user").(*userModel.User)
	if !exists {
		return response.Unauthorized(c, "unauthorized", nil)
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid category id", err)
	}

	var input categoryModel.CreateCategoryInput
	if err := c.Bind(&input); err != nil {
		return response.BadRequest(c, "invalid request payload", err)
	}

	category, err := h.categoryService.Update(uint(id), input, user)
	if err != nil {
		return response.BadRequest(c, "failed to update category", err)
	}

	return response.Success(c, http.StatusOK, "Category updated successfully", category)
}

func (h *CategoryHandler) Delete(c echo.Context) error {
	user, exists := c.Get("user").(*userModel.User)
	if !exists {
		return response.Unauthorized(c, "unauthorized", nil)
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return response.BadRequest(c, "invalid category id", err)
	}

	if err := h.categoryService.Delete(uint(id), user); err != nil {
		return response.BadRequest(c, "failed to delete category", err)
	}

	return response.Success(c, http.StatusOK, "Category deleted successfully", nil)
}
