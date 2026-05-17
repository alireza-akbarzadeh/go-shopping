package services

import (
	"context"
	"errors"
	"strings"

	"github.com/alireza-akbarzadeh/luxe/internal/dto"
	"github.com/alireza-akbarzadeh/luxe/internal/models"
	"github.com/alireza-akbarzadeh/luxe/internal/utils"
	"gorm.io/gorm"
)

type UserMenuServicesInterface interface {
	GetAllGroups() ([]models.MenuGroup, error)
	GetGroupByID(id uint) (*models.MenuGroup, error)
	CreateGroup(req *dto.CreateMenuGroupRequest) (*models.MenuGroup, error)
	UpdateGroup(id uint, req *dto.UpdateMenuGroupRequest) (*models.MenuGroup, error)
	DeleteGroup(id uint) error

	// Items
	GetAllItems(flat bool) ([]models.MenuItem, error)
	GetItemByID(id uint) (*models.MenuItem, error)
	CreateItem(req *dto.CreateMenuItemRequest) (*models.MenuItem, error)
	UpdateItem(id uint, req *dto.UpdateMenuItemRequest) (*models.MenuItem, error)
	DeleteItem(id uint) error

	// User-facing: returns filtered sidebar for given role and search term
	GetUserMenu(ctx context.Context, userRole string, search string) ([]dto.SidebarGroup, error)
}

type userMenuService struct {
	db *gorm.DB
}

func NewMenuService(db *gorm.DB) UserMenuServicesInterface {
	return &userMenuService{db: db}
}

// GetAllGroups menus
func (s *userMenuService) GetAllGroups() ([]models.MenuGroup, error) {
	var groups []models.MenuGroup
	err := s.db.Order("display_order ASC").Find(&groups).Error
	return groups, err
}

// find the gorup ids
func (s *userMenuService) GetGroupByID(id uint) (*models.MenuGroup, error) {
	var group models.MenuGroup
	err := s.db.First(&group, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &group, err
}

func (s *userMenuService) CreateGroup(req *dto.CreateMenuGroupRequest) (*models.MenuGroup, error) {
	group := &models.MenuGroup{
		Name:         req.Name,
		DisplayOrder: req.DisplayOrder,
	}
	err := s.db.Create(group).Error
	return group, err
}

func (s *userMenuService) UpdateGroup(id uint, req *dto.UpdateMenuGroupRequest) (*models.MenuGroup, error) {
	group, err := s.GetGroupByID(id)

	if err != nil || group == nil {
		return nil, utils.ErrNotFound("group not found")
	}
	group.Name = req.Name
	group.DisplayOrder = req.DisplayOrder
	err = s.db.Save(group).Error
	return group, err
}

func (s *userMenuService) DeleteGroup(id uint) error {
	return s.db.Delete(&models.MenuGroup{}, id).Error
}

func (s *userMenuService) GetAllItems(flat bool) ([]models.MenuItem, error) {
	var allItems []models.MenuItem
	err := s.db.Order("display_order ASC").Find(&allItems).Error
	if err != nil {
		return nil, err
	}

	if flat {
		return allItems, nil
	}

	return s.buildTree(allItems), nil
}

func (s *userMenuService) buildTree(items []models.MenuItem) []models.MenuItem {
	itemMap := make(map[uint]*models.MenuItem)
	for i := range items {
		itemMap[items[i].ID] = &items[i]
	}

	var roots []models.MenuItem
	for i := range items {
		parentID := items[i].ParentID
		if parentID == nil || *parentID == 0 {
			roots = append(roots, items[i])
		} else {
			if parent, exists := itemMap[*parentID]; exists {
				parent.Children = append(parent.Children, items[i])
			}
		}
	}
	return roots
}

func (s *userMenuService) GetItemByID(id uint) (*models.MenuItem, error) {
	var item models.MenuItem
	err := s.db.First(&item, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &item, err
}

func (s *userMenuService) CreateItem(req *dto.CreateMenuItemRequest) (*models.MenuItem, error) {
	var group models.MenuGroup
	if err := s.db.First(&group, req.GroupID).Error; err != nil {
		return nil, errors.New("invalid group_id")
	}
	if req.ParentID != nil && *req.ParentID != 0 {
		var parent models.MenuItem
		if err := s.db.First(&parent, *req.ParentID).Error; err != nil {
			return nil, utils.ErrBadRequest("invalid parent_id")
		}
		if parent.GroupID != req.GroupID {
			return nil, utils.ErrBadRequest("parent_id must belong to the same group")
		}
	}
	item := &models.MenuItem{
		GroupID:      req.GroupID,
		ParentID:     req.ParentID,
		Label:        req.Label,
		Href:         req.Href,
		Icon:         req.Icon,
		Permission:   req.Permission,
		DisplayOrder: req.DisplayOrder,
	}
	err := s.db.Create(item).Error
	return item, err
}

func (s *userMenuService) UpdateItem(id uint, req *dto.UpdateMenuItemRequest) (*models.MenuItem, error) {
	item, err := s.GetItemByID(id)
	if err != nil || item == nil {
		return nil, errors.New("item not found")
	}
	// validate group if changed
	if item.GroupID != req.GroupID {
		var group models.MenuGroup
		if err := s.db.First(&group, req.GroupID).Error; err != nil {
			return nil, errors.New("invalid group_id")
		}
		item.GroupID = req.GroupID
	}
	if req.ParentID != nil && *req.ParentID != 0 {
		var parent models.MenuItem
		if err := s.db.First(&parent, *req.ParentID).Error; err != nil {
			return nil, errors.New("invalid parent_id")
		}
		if parent.GroupID != req.GroupID {
			return nil, errors.New("parent_id must belong to same group")
		}
		item.ParentID = req.ParentID
	} else {
		item.ParentID = nil
	}
	item.Label = req.Label
	item.Href = req.Href
	item.Icon = req.Icon
	item.Permission = req.Permission
	item.DisplayOrder = req.DisplayOrder
	err = s.db.Save(item).Error
	return item, err
}

func (s *userMenuService) DeleteItem(id uint) error {
	return s.db.Delete(&models.MenuItem{}, id).Error
}

func (s *userMenuService) GetUserMenu(ctx context.Context, userRole string, search string) ([]dto.SidebarGroup, error) {
	var groups []models.MenuGroup
	err := s.db.Preload("Items", func(db *gorm.DB) *gorm.DB {
		return db.Order("display_order ASC")
	}).Order("display_order ASC").Find(&groups).Error
	if err != nil {
		return nil, err
	}

	var allItems []models.MenuItem
	for _, g := range groups {
		allItems = append(allItems, g.Items...)
	}

	itemMap, roots := s.buildFullTree(allItems)

	var result []dto.SidebarGroup
	for _, group := range groups {
		var groupRoots []models.MenuItem
		for _, root := range roots {
			if root.GroupID == group.ID {
				groupRoots = append(groupRoots, root)
			}
		}
		visibleItems := s.filterItems(groupRoots, userRole, search, itemMap)
		if len(visibleItems) > 0 {
			result = append(result, dto.SidebarGroup{
				Group: group.Name,
				Items: visibleItems,
			})
		}
	}
	return result, nil
}

func (s *userMenuService) filterItems(items []models.MenuItem, userRole, search string, itemMap map[uint]*models.MenuItem) []dto.SidebarItem {
	var result []dto.SidebarItem
	for _, item := range items {
		// Role check
		if item.Permission != nil && *item.Permission != userRole && userRole != "admin" {
			continue
		}
		// Recursively filter children
		filteredChildren := s.filterItems(item.Children, userRole, search, itemMap)
		// Search match on current item
		matchesSearch := search == "" || s.matchesSearch(item.Label, item.Href, search)
		if search != "" && !matchesSearch && len(filteredChildren) == 0 {
			continue
		}
		sidebarItem := dto.SidebarItem{
			Label:    item.Label,
			Icon:     item.Icon,
			Children: filteredChildren,
		}
		if item.Href != nil {
			sidebarItem.Href = *item.Href
		}
		result = append(result, sidebarItem)
	}
	return result
}

func (s *userMenuService) buildFullTree(items []models.MenuItem) (map[uint]*models.MenuItem, []models.MenuItem) {
	itemMap := make(map[uint]*models.MenuItem)
	for i := range items {
		itemMap[items[i].ID] = &items[i]
		itemMap[items[i].ID].Children = []models.MenuItem{} // initialize
	}
	var roots []models.MenuItem
	for i := range items {
		parentID := items[i].ParentID
		if parentID == nil || *parentID == 0 {
			roots = append(roots, items[i])
		} else {
			if parent, ok := itemMap[*parentID]; ok {
				parent.Children = append(parent.Children, items[i])
			}
		}
	}
	return itemMap, roots
}

func (s *userMenuService) matchesSearch(label string, href *string, search string) bool {
	searchLower := strings.ToLower(search)
	if strings.Contains(strings.ToLower(label), searchLower) {
		return true
	}
	if href != nil && strings.Contains(strings.ToLower(*href), searchLower) {
		return true
	}
	return false
}
