package model

import (
	"errors"
	"time"
)

type Category struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedBy   uint      `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CreateCategoryInput struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func NewCategory(input CreateCategoryInput, userID uint) (*Category, error) {
	if input.Name == "" {
		return nil, errors.New("category name is required")
	}

	return &Category{
		Name:        input.Name,
		Description: input.Description,
		CreatedBy:   userID,
	}, nil
}

func (c *Category) Update(input CreateCategoryInput) error {
	if input.Name == "" {
		return errors.New("category name is required")
	}

	c.Name = input.Name
	c.Description = input.Description
	return nil
}
