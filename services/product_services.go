package services

import (
	"errors"
	"fmt"

	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"gorm.io/gorm"
)

type ProductServiceInterface interface {
	List(limit, offset int, filters map[string]interface{}) ([]models.Product, int64, error)
	GetByID(id uint) (*models.Product, error)
	GetBySlug(slug string) (*models.Product, error)
	Create(req CreateProductRequest) (*models.Product, error)
	Update(productID uint, req UpdateProductRequest) (*models.Product, error)
	Delete(id uint) error
}

type CreateProductRequest struct {
	Name              string   `json:"name" validate:"required,min=3,max=255"`
	Description       string   `json:"description,omitempty"`
	Price             float64  `json:"price" validate:"required,gte=0"`
	CompareAtPrice    *float64 `json:"compare_at_price,omitempty" validate:"omitempty,gte=0"`
	Cost              *float64 `json:"cost,omitempty" validate:"omitempty,gte=0"`
	SKU               string   `json:"sku" validate:"required,min=3,max=50"`
	Barcode           string   `json:"barcode,omitempty"`
	Stock             int      `json:"stock" validate:"gte=0"`
	LowStockThreshold int      `json:"low_stock_threshold,omitempty"`
	Weight            *float64 `json:"weight,omitempty" validate:"omitempty,gte=0"`
	IsDigital         bool     `json:"is_digital"`
	CategoryID        *uint    `json:"category_id,omitempty"`
	Images            []string `json:"images,omitempty"`
	Status            string   `json:"status" validate:"oneof=draft active inactive archived"`
	MetaTitle         string   `json:"meta_title,omitempty"`
	MetaDescription   string   `json:"meta_description,omitempty"`
}
type UpdateProductRequest struct {
	Name              *string   `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description       *string   `json:"description,omitempty"`
	Price             *float64  `json:"price,omitempty" validate:"omitempty,gte=0"`
	CompareAtPrice    *float64  `json:"compare_at_price,omitempty" validate:"omitempty,gte=0"`
	Cost              *float64  `json:"cost,omitempty" validate:"omitempty,gte=0"`
	SKU               *string   `json:"sku,omitempty" validate:"omitempty,min=3,max=50"`
	Barcode           *string   `json:"barcode,omitempty"`
	Stock             *int      `json:"stock,omitempty" validate:"omitempty,gte=0"`
	LowStockThreshold *int      `json:"low_stock_threshold,omitempty"`
	Weight            *float64  `json:"weight,omitempty" validate:"omitempty,gte=0"`
	IsDigital         *bool     `json:"is_digital,omitempty"`
	CategoryID        *uint     `json:"category_id,omitempty"`
	Images            *[]string `json:"images,omitempty"`
	Status            *string   `json:"status,omitempty" validate:"omitempty,oneof=draft active inactive archived"`
	MetaTitle         *string   `json:"meta_title,omitempty"`
	MetaDescription   *string   `json:"meta_description,omitempty"`
}

type productService struct {
	db *gorm.DB
}

func NewProductService(db *gorm.DB) ProductServiceInterface {
	return &productService{db: db}
}

// ensureUniqueSlug checks and modifies slug to be unique.
func (s *productService) UniqSlug(baseSlug string, excludeID uint) string {
	slug := baseSlug
	counter := 1
	for {
		var count int64
		query := s.db.Model(&models.Product{}).Where("slug = ?", slug)
		if excludeID > 0 {
			query = query.Where("id != ?", excludeID)
		}
		query.Count(&count)
		if count == 0 {
			break
		}
		slug = fmt.Sprintf("%s-%d", baseSlug, counter)
		counter++
	}
	return slug
}

func (s *productService) Create(req CreateProductRequest) (*models.Product, error) {
	baseSlug := generateSlug(req.Name)
	slug := s.UniqSlug(baseSlug, 0)
	product := models.Product{
		Name:              req.Name,
		Slug:              slug,
		Description:       req.Description,
		Price:             req.Price,
		CompareAtPrice:    req.CompareAtPrice,
		Cost:              req.Cost,
		SKU:               req.SKU,
		Barcode:           req.Barcode,
		Stock:             req.Stock,
		LowStockThreshold: req.LowStockThreshold,
		Weight:            req.Weight,
		IsDigital:         req.IsDigital,
		CategoryID:        req.CategoryID,
		Images:            req.Images,
		Status:            req.Status,
		MetaTitle:         req.MetaTitle,
		MetaDescription:   req.MetaDescription,
	}
	if product.Status == "" {
		product.Status = "draft"
	}
	if product.LowStockThreshold == 0 {
		product.LowStockThreshold = 5
	}

	if err := s.db.Create(product).Error; err != nil {
		return nil, utils.ErrInternal(err)
	}
	return &product, nil
}

// Retrive product by id
func (s *productService) GetByID(id uint) (*models.Product, error) {
	var product models.Product
	if err := s.db.Preload("Category").First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("product not found")
		}
		return nil, utils.ErrInternal(err)
	}
	return &product, nil
}

// Retrive product by slug
func (s *productService) GetBySlug(slug string) (*models.Product, error) {
	var product models.Product
	if err := s.db.Preload("Category").Where("slug = ?", slug).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("product not found")
		}
		return nil, utils.ErrInternal(err)
	}
	return &product, nil
}

// Update product

func (s *productService) Update(id uint, req UpdateProductRequest) (*models.Product, error) {
	product, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		product.Name = *req.Name
		baseSlug := generateSlug(*req.Name)
		product.Slug = s.UniqSlug(baseSlug, id)
	}
	if req.Description != nil {
		product.Description = *req.Description
	}
	if req.Price != nil {
		product.Price = *req.Price
	}
	if req.CompareAtPrice != nil {
		product.CompareAtPrice = req.CompareAtPrice
	}
	if req.Cost != nil {
		product.Cost = req.Cost
	}
	if req.SKU != nil {
		var existing models.Product
		if err := s.db.Where("sku = ? AND id != ?", *req.SKU, id).First(&existing).Error; err == nil {
			return nil, utils.ErrConflict("SKU already exists")
		}
		product.SKU = *req.SKU
	}
	if req.Barcode != nil {
		product.Barcode = *req.Barcode
	}
	if req.Stock != nil {
		product.Stock = *req.Stock
	}
	if req.LowStockThreshold != nil {
		product.LowStockThreshold = *req.LowStockThreshold
	}
	if req.Weight != nil {
		product.Weight = req.Weight
	}
	if req.IsDigital != nil {
		product.IsDigital = *req.IsDigital
	}
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}
	if req.Images != nil {
		product.Images = *req.Images
	}
	if req.Status != nil {
		product.Status = *req.Status
	}
	if req.MetaTitle != nil {
		product.MetaTitle = *req.MetaTitle
	}
	if req.MetaDescription != nil {
		product.MetaDescription = *req.MetaDescription
	}
	if err := s.db.Save(product).Error; err != nil {
		return nil, utils.ErrInternal(err)
	}
	return product, nil
}

// Delete product
func (s *productService) Delete(id uint) error {
	result := s.db.Delete(&models.Product{}, id)
	if result.Error != nil {
		return utils.ErrInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.ErrNotFound("product not found")
	}
	return nil
}

// retrieve list of product
func (s *productService) List(limit, offset int, filters map[string]interface{}) ([]models.Product, int64, error) {
	const maxLimit = 100
	if limit <= 0 {
		limit = 10
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	if offset < 0 {
		offset = 0
	}

	query := s.db.Model(&models.Product{}).Order("id DESC")

	// String filters
	if v, ok := filters["status"].(string); ok && v != "" {
		query = query.Where("status = ?", v)
	}
	if v, ok := filters["name"].(string); ok && v != "" {
		query = query.Where("name = ?", v)
	}
	if v, ok := filters["sku"].(string); ok && v != "" {
		query = query.Where("sku = ?", v)
	}

	// Numeric filters
	if v, ok := filters["category_id"].(uint); ok && v != 0 {
		query = query.Where("category_id = ?", v)
	}
	if v, ok := filters["min_price"].(float64); ok {
		query = query.Where("price >= ?", v)
	}
	if v, ok := filters["max_price"].(float64); ok {
		query = query.Where("price <= ?", v)
	}

	// Boolean filter
	if v, ok := filters["is_digital"].(bool); ok {
		query = query.Where("is_digital = ?", v)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count products: %w", err)
	}

	var products []models.Product
	if err := query.Limit(limit).Offset(offset).Preload("Category").Find(&products).Error; err != nil {
		return nil, 0, fmt.Errorf("find products: %w", err)
	}

	return products, total, nil
}
