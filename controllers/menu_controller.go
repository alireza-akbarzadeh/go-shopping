package controllers

import (
	"net/http"
	"strconv"

	"github.com/alireza-akbarzadeh/shopping-platform/dto"
	"github.com/alireza-akbarzadeh/shopping-platform/services"
	"github.com/gin-gonic/gin"
)

type MenuController struct {
	menuService services.UserMenuServicesInterface
}

func NewMenuController(menuService services.UserMenuServicesInterface) *MenuController {
	return &MenuController{menuService: menuService}
}

// GetAllGroups godoc
// @Summary      Get all menu groups
// @Description  Retrieves all menu groups ordered by display_order
// @Tags         Admin Menu Groups
// @Produce      json
// @Success      200  {array}   models.MenuGroup
// @Failure      500  {object}  dto.MessageResponse
// @Router       /admin/menu/groups [get]
// @Security     BearerAuth

func (ctrl *MenuController) GetAllGroups(c *gin.Context) {
	groups, err := ctrl.menuService.GetAllGroups()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, groups)
}

// GetGroupByID godoc
// @Summary      Get group by ID
// @Description  Returns a single menu group by its ID
// @Tags         Admin Menu Groups
// @Produce      json
// @Param        id   path      int  true  "Group ID"
// @Success      200  {object}  models.MenuGroup
// @Failure      400  {object}  dto.MessageResponse
// @Failure      404  {object}  dto.MessageResponse
// @Failure      500  {object}  dto.MessageResponse
// @Router       /admin/menu/groups/{id} [get]
// @Security     BearerAuth
func (ctrl *MenuController) GetGroupByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}
	group, err := ctrl.menuService.GetGroupByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if group == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "group not found"})
		return
	}
	c.JSON(http.StatusOK, group)
}

// CreateGroup godoc
// @Summary      Create a new menu group
// @Description  Creates a menu group (e.g., "Overview", "Users & Access")
// @Tags         Admin Menu Groups
// @Accept       json
// @Produce      json
// @Param        request body      dto.CreateMenuGroupRequest true "Group data"
// @Success      201     {object}  models.MenuGroup
// @Failure      400     {object}  dto.MessageResponse
// @Failure      500     {object}  dto.MessageResponse
// @Router       /admin/menu/groups [post]
// @Security     BearerAuth
func (ctrl *MenuController) CreateGroup(c *gin.Context) {
	var req dto.CreateMenuGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	group, err := ctrl.menuService.CreateGroup(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, group)
}

// UpdateGroup godoc
// @Summary      Update an existing menu group
// @Description  Updates group name or display order
// @Tags         Admin Menu Groups
// @Accept       json
// @Produce      json
// @Param        id       path      int                        true "Group ID"
// @Param        request  body      dto.UpdateMenuGroupRequest true "Updated group data"
// @Success      200      {object}  models.MenuGroup
// @Failure      400      {object}  dto.MessageResponse
// @Failure      404      {object}  dto.MessageResponse
// @Failure      500      {object}  dto.MessageResponse
// @Router       /admin/menu/groups/{id} [put]
// @Security     BearerAuth
func (ctrl *MenuController) UpdateGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}
	var req dto.UpdateMenuGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	group, err := ctrl.menuService.UpdateGroup(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, group)
}

// DeleteGroup godoc
// @Summary      Delete a menu group
// @Description  Deletes a group and all its menu items (cascade)
// @Tags         Admin Menu Groups
// @Produce      json
// @Param        id   path      int  true "Group ID"
// @Success      204  "No Content"
// @Failure      400  {object}  dto.MessageResponse
// @Failure      500  {object}  dto.MessageResponse
// @Router       /admin/menu/groups/{id} [delete]
// @Security     BearerAuth
func (ctrl *MenuController) DeleteGroup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid group id"})
		return
	}
	if err := ctrl.menuService.DeleteGroup(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetAllItems godoc
// @Summary      Get all menu items
// @Description  Returns flat or nested menu items (use ?flat=true for flat list)
// @Tags         Admin Menu Items
// @Produce      json
// @Param        flat  query   bool  false  "Return flat list" default(false)
// @Success      200   {array}  models.MenuItem
// @Failure      500   {object}  dto.MessageResponse
// @Router       /admin/menu/items [get]
// @Security     BearerAuth
func (ctrl *MenuController) GetAllItems(c *gin.Context) {
	if ctrl.menuService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "menu service not initialized"})
		return
	}
	flat, _ := strconv.ParseBool(c.DefaultQuery("flat", "false"))
	items, err := ctrl.menuService.GetAllItems(flat)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// GetItemByID godoc
// @Summary      Get menu item by ID
// @Description  Returns a single menu item by its ID
// @Tags         Admin Menu Items
// @Produce      json
// @Param        id   path      int  true "Item ID"
// @Success      200  {object}  models.MenuItem
// @Failure      400  {object}  dto.MessageResponse
// @Failure      404  {object}  dto.MessageResponse
// @Failure      500  {object}  dto.MessageResponse
// @Router       /admin/menu/items/{id} [get]
// @Security     BearerAuth
func (ctrl *MenuController) GetItemByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}
	item, err := ctrl.menuService.GetItemByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if item == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}
	c.JSON(http.StatusOK, item)
}

// CreateItem godoc
// @Summary      Create a new menu item
// @Description  Adds a new menu item (can be top-level or child of another item)
// @Tags         Admin Menu Items
// @Accept       json
// @Produce      json
// @Param        request body      dto.CreateMenuItemRequest true "Menu item data"
// @Success      201     {object}  models.MenuItem
// @Failure      400     {object}  dto.MessageResponse
// @Failure      500     {object}  dto.MessageResponse
// @Router       /admin/menu/items [post]
// @Security     BearerAuth
func (ctrl *MenuController) CreateItem(c *gin.Context) {
	var req dto.CreateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := ctrl.menuService.CreateItem(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

// UpdateItem godoc
// @Summary      Update an existing menu item
// @Description  Updates menu item details including group, parent, label, href, etc.
// @Tags         Admin Menu Items
// @Accept       json
// @Produce      json
// @Param        id       path      int                       true "Item ID"
// @Param        request  body      dto.UpdateMenuItemRequest true "Updated item data"
// @Success      200      {object}  models.MenuItem
// @Failure      400      {object}  dto.MessageResponse
// @Failure      404      {object}  dto.MessageResponse
// @Failure      500      {object}  dto.MessageResponse
// @Router       /admin/menu/items/{id} [put]
// @Security     BearerAuth
func (ctrl *MenuController) UpdateItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}
	var req dto.UpdateMenuItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := ctrl.menuService.UpdateItem(uint(id), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, item)
}

// DeleteItem godoc
// @Summary      Delete a menu item
// @Description  Deletes a menu item and all its children (cascade delete due to foreign key constraint)
// @Tags         Admin Menu Items
// @Produce      json
// @Param        id   path      int  true  "Menu item ID"
// @Success      204  "No Content"
// @Failure      400  {object}  dto.MessageResponse
// @Failure      500  {object}  dto.MessageResponse
// @Router       /admin/menu/items/{id} [delete]
// @Security     BearerAuth
func (ctrl *MenuController) DeleteItem(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}
	if err := ctrl.menuService.DeleteItem(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}

// GetUserMenu godoc
// @Summary      Get user sidebar menu
// @Description  Returns the sidebar menu filtered by user's role and optional search term
// @Tags         User Menu
// @Produce      json
// @Param        search  query   string  false  "Search by label or href"
// @Success      200     {array}  dto.SidebarGroup
// @Failure      500     {object}  dto.MessageResponse
// @Router       /user/menu [get]
// @Security     BearerAuth
func (ctrl *MenuController) GetUserMenu(c *gin.Context) {
	// Extract user role from context (set by auth middleware)
	userRole, exists := c.Get("user_role")
	if !exists {
		userRole = "guest" // default role
	}
	search := c.Query("search")
	menu, err := ctrl.menuService.GetUserMenu(c.Request.Context(), userRole.(string), search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, menu)
}
