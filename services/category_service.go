package services

import (
	"errors"
	"fmt"

	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"gorm.io/gorm"
)

type CategoryServiceInterface interface {
	Create(req CreateCategoryRequest) (*models.Category, error)
	GetByID(id uint) (*models.Category, error)
	GetBySlug(slug string) (*models.Category, error)
	Update(id uint, req UpdateCategoryRequest) (*models.Category, error)
	Delete(id uint) error
	List(limit, offset int, filters map[string]interface{}) ([]models.Category, int64, error)
	BulkCreate(categories []CreateCategoryRequest) ([]*models.Category, error)
	BulkDelete(ids []uint) error
}

type categoryService struct {
	db *gorm.DB
}

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Slug        string `json:"slug" validate:"required,slug"`
	Description string `json:"description,omitempty"`
	ParentID    *uint  `json:"parent_id,omitempty"`
	IsActive    bool   `json:"is_active"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Slug        *string `json:"slug,omitempty" validate:"omitempty,slug"`
	Description *string `json:"description,omitempty"`
	ParentID    *uint   `json:"parent_id,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

func NewCategoryService(db *gorm.DB) CategoryServiceInterface {
	return &categoryService{db: db}
}

func (s *categoryService) uniqueCategorySlug(baseSlug string, excludeID uint) string {
	slug := baseSlug
	counter := 1
	for {
		var count int64
		query := s.db.Model(&models.Category{}).Where("slug = ?", slug)
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

// updateLevelAndPath recalculates level and path based on parent
func (s *categoryService) updateLevelAndPath(category *models.Category) error {
	if category.ParentID == nil || *category.ParentID == 0 {
		category.Level = 0
		category.Path = ""
		return nil
	}
	var parent models.Category
	if err := s.db.First(&parent, *category.ParentID).Error; err != nil {
		return err
	}
	category.Level = parent.Level + 1
	if parent.Path == "" {
		category.Path = fmt.Sprintf("%d", parent.ID)
	} else {
		category.Path = fmt.Sprintf("%s.%d", parent.Path, parent.ID)
	}
	return nil

}

// Create categories
func (s *categoryService) Create(req CreateCategoryRequest) (*models.Category, error) {
	slug := req.Slug
	if slug == "" {
		slug = generateSlug(req.Name)
	}
	slug = s.uniqueCategorySlug(slug, uint(0))

	category := &models.Category{
		Name:        slug,
		Slug:        slug,
		Description: req.Description,
		ParentID:    req.ParentID,
		IsActive:    req.IsActive,
	}
	if err := s.updateLevelAndPath(category); err != nil {
		return nil, utils.ErrInternal(err)
	}
	if err := s.db.Create(category).Error; err != nil {
		return nil, utils.ErrInternal(err)
	}
	return category, nil
}

// GetByID retrieve category by id
func (s *categoryService) GetByID(id uint) (*models.Category, error) {
	var category models.Category
	if err := s.db.Preload("Parent").Preload("Children").First(&category, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("category not found")
		}
		return nil, utils.ErrInternal(err)
	}
	return &category, nil
}

// GetBySlug retrieve category by slug
func (s *categoryService) GetBySlug(slug string) (*models.Category, error) {
	var category models.Category
	if err := s.db.Preload("Parent").Preload("Children").Where("slug = ?", slug).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("category not found")
		}
		return nil, utils.ErrInternal(err)
	}
	return &category, nil
}

// Update update categories
func (s *categoryService) Update(id uint, req UpdateCategoryRequest) (*models.Category, error) {
	category, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		category.Name = *req.Name
		if req.Slug == nil || *req.Slug == "" {
			baseSlug := generateSlug(*req.Name)
			category.Slug = s.uniqueCategorySlug(baseSlug, id)
		}
	}
	if req.Slug != nil && *req.Slug != "" {
		category.Slug = s.uniqueCategorySlug(*req.Slug, id)
	}
	if req.Description != nil {
		category.Description = *req.Description
	}
	if req.ParentID != nil {
		category.ParentID = req.ParentID
		if err := s.updateLevelAndPath(category); err != nil {
			return nil, utils.ErrInternal(err)
		}
	}
	if req.IsActive != nil {
		category.IsActive = *req.IsActive
	}

	if err := s.db.Save(category).Error; err != nil {
		return nil, utils.ErrInternal(err)
	}
	return category, nil
}

// Delete remove categories
func (s *categoryService) Delete(id uint) error {
	var count int64
	s.db.Model(&models.Category{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		return utils.ErrBadRequest("cannot delete category with children; delete children first")
	}
	result := s.db.Delete(&models.Category{}, id)
	if result.Error != nil {
		return utils.ErrInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.ErrNotFound("category not found")
	}
	return nil
}

// List retrieve list for categories
func (s *categoryService) List(limit, offset int, filters map[string]interface{}) ([]models.Category, int64, error) {
	var categories []models.Category
	var total int64

	query := s.db.Model(&models.Category{})

	if isActive, ok := filters["is_active"]; ok && isActive != nil {
		query = query.Where("is_active = ?", isActive)
	}
	if parentID, ok := filters["parent_id"]; ok && parentID != nil {
		query = query.Where("parent_id = ?", parentID)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, utils.ErrInternal(err)
	}
	if err := query.Limit(limit).Offset(offset).Preload("Parent").Preload("Children").Find(&categories).Error; err != nil {
		return nil, 0, utils.ErrInternal(err)
	}
	return categories, total, nil
}

//BulkCreate

func (s *categoryService) BulkCreate(categories []CreateCategoryRequest) ([]*models.Category, error) {
	if len(categories) == 0 {
		return nil, utils.ErrBadRequest("no categories provided")
	}
	var created []*models.Category
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, req := range categories {
			slug := req.Slug
			if slug == "" {
				slug = generateSlug(req.Name)
			}
			slug = s.uniqueCategorySlug(slug, 0)

			cat := &models.Category{
				Name:        req.Name,
				Slug:        slug,
				Description: req.Description,
				ParentID:    req.ParentID,
				IsActive:    req.IsActive,
			}
			if err := s.updateLevelAndPath(cat); err != nil {
				return err
			}
			if err := tx.Create(cat).Error; err != nil {
				return err
			}
			created = append(created, cat)
		}
		return nil
	})
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	return created, nil
}

// BulkDelete remove categorie with given ids
func (s *categoryService) BulkDelete(ids []uint) error {
	if len(ids) == 0 {
		return utils.ErrBadRequest("no category IDs provided")
	}

	// Check if any category has children
	var count int64
	s.db.Model(&models.Category{}).Where("parent_id IN ?", ids).Count(&count)
	if count > 0 {
		return utils.ErrBadRequest("cannot delete categories that have children; delete children first")
	}

	result := s.db.Where("id IN ?", ids).Delete(&models.Category{})
	if result.Error != nil {
		return utils.ErrInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.ErrNotFound("no categories found to delete")
	}
	return nil
}
