package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Quantity    int            `gorm:"default:0;not null" json:"quantity"`
	Status      int16          `gorm:"default:1" json:"status"`
	ImageURL    string         `gorm:"type:varchar(255)" json:"image_url"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type CreateProductInput struct {
	Name        string  `form:"name" binding:"required,not_blank,min=3"`
	Description string  `form:"description" binding:"required,not_blank"`
	Price       float64 `form:"price" binding:"required,gt=0"`
	Quantity    int     `form:"quantity" binding:"required,gte=0"`
}

type UpdateProductInput struct {
	Name        string  `form:"name" binding:"required,not_blank,min=3"`
	Description string  `form:"description" binding:"required,not_blank"`
	Price       float64 `form:"price" binding:"required,gt=0"`
	Quantity    int     `form:"quantity" binding:"required,gte=0"` // ✅ Diperbaiki
	ImageURL    string  `form:"image_url"`
}
