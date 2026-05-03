package controllers

import (
	"strconv"

	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/services"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CategoryController struct {
	categoryService services.CategoryServiceInterface
	validate        *validator.Validate
}

func NewCategoryController(categoryService services.CategoryServiceInterface) *CategoryController {
	return &CategoryController{
		categoryService: categoryService,
		validate:        validator.New(),
	}
}

// Create creates a new category (admin only).
// @Summary      Create a new category
// @Description  Creates a new product category. Only accessible by users with the "admin" role.
// @Tags         Admin Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body services.CreateCategoryRequest true "Category creation data"
// @Success      201 {object} utils.Response{data=models.Category}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      403 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /categories [post]
func (ctrl *CategoryController) Create(c *gin.Context) {
	var req services.CreateCategoryRequest
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}
	category, err := ctrl.categoryService.Create(req)
	if err != nil {
		utils.InternalServerErrorResponse(c, err, "failed to create category")
		return
	}
	utils.CreatedResponse(c, constants.MsgCreateSuccess, category)

}

// Update updates an existing category (admin only).
// @Summary      Update a category
// @Description  Updates an existing category by ID. Only accessible by users with the "admin" role.
// @Tags         Admin Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path      int                           true  "Category ID"
// @Param        request body      services.UpdateCategoryRequest true  "Category update data"
// @Success      200     {object}  utils.Response{data=models.Category}
// @Failure      400     {object}  utils.Response
// @Failure      401     {object}  utils.Response
// @Failure      403     {object}  utils.Response
// @Failure      404     {object}  utils.Response
// @Failure      500     {object}  utils.Response
// @Router       /categories/{id} [put]
func (ctrl *CategoryController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid category id")
		return
	}

	var req services.UpdateCategoryRequest
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}

	category, err := ctrl.categoryService.Update(uint(id), req)
	if err != nil {
		utils.HandleAppError(c, err, "failed to update category")
		return
	}
	utils.SuccessResponse(c, constants.MsgUpdateSuccess, category)
}

// Delete deletes a category (admin only).
// @Summary      Delete a category
// @Description  Deletes a category by ID. Only accessible by users with the "admin" role.
// @Description  Categories with child categories cannot be deleted – delete children first.
// @Tags         Admin Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  utils.Response
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /categories/{id} [delete]
func (ctrl *CategoryController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid category id")
		return
	}

	if err := ctrl.categoryService.Delete(uint(id)); err != nil {
		utils.HandleAppError(c, err, "failed to delete category")
		return
	}
	utils.SuccessResponse(c, constants.MsgDeleteSuccess, nil)
}

// GetOne retrieves a single category by ID or slug (public).
// @Summary      Get a category by ID or slug
// @Description  Returns a single category. Accepts either a numeric ID or a URL slug.
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        identifier   path      string  true  "Category ID (numeric) or slug (string)"
// @Success      200          {object}  utils.Response{data=models.Category}
// @Failure      400          {object}  utils.Response
// @Failure      404          {object}  utils.Response
// @Failure      500          {object}  utils.Response
// @Router       /categories/{identifier} [get]
func (ctrl *CategoryController) GetOne(c *gin.Context) {
	identifier := c.Param("identifier")
	id, err := strconv.ParseUint(identifier, 10, 64)
	var category *models.Category
	if err == nil {
		category, err = ctrl.categoryService.GetByID(uint(id))
	} else {
		category, err = ctrl.categoryService.GetBySlug(identifier)
	}
	if err != nil {
		utils.HandleAppError(c, err, "failed to fetch category")
		return
	}
	utils.SuccessResponse(c, constants.MsgFetchSuccess, category)
}

// List returns a paginated list of categories (public).
// @Summary      List categories
// @Description  Returns a paginated list of categories with optional filtering.
// @Tags         Categories
// @Accept       json
// @Produce      json
// @Param        limit       query     int     false  "Items per page"  default(20)  minimum(1)  maximum(100)
// @Param        offset      query     int     false  "Offset (skip number of items)"  default(0)  minimum(0)
// @Param        is_active   query     bool    false  "Filter by active status (true/false)"
// @Param        parent_id   query     int     false  "Filter by parent category ID"
// @Success      200         {object}  utils.Response{data=object{categories=[]models.Category,total=int,limit=int,offset=int}}
// @Failure      400         {object}  utils.Response
// @Failure      500         {object}  utils.Response
// @Router       /categories [get]
func (ctrl *CategoryController) List(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	filters := make(map[string]interface{})
	if isActive := c.Query("is_active"); isActive != "" {
		filters["is_active"] = isActive == "true"
	}
	if parentID := c.Query("parent_id"); parentID != "" {
		if pid, err := strconv.ParseUint(parentID, 10, 64); err == nil {
			filters["parent_id"] = uint(pid)
		}
	}

	categories, total, err := ctrl.categoryService.List(limit, offset, filters)
	if err != nil {
		utils.InternalServerErrorResponse(c, err, "failed to list categories")
		return
	}
	data := gin.H{
		"categories": categories,
		"total":      total,
		"limit":      limit,
		"offset":     offset,
	}
	utils.SuccessResponse(c, constants.MsgFetchSuccess, data)
}

// BulkCreate creates multiple categories in one request (admin only).
// @Summary      Bulk create categories
// @Description  Creates multiple categories at once. Only accessible by users with the "admin" role.
// @Tags         Admin Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body []services.CreateCategoryRequest true "Array of categories to create"
// @Success      201 {object} utils.Response{data=[]models.Category}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      403 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /categories/bulk [post]
func (ctrl *CategoryController) BulkCreate(c *gin.Context) {
	var reqs []services.CreateCategoryRequest
	if !utils.BindAndValidate(c, &reqs, ctrl.validate) {
		return
	}
	if len(reqs) == 0 {
		utils.ErrorResponse(c, 400, "no categories provided")
		return
	}
	categories, err := ctrl.categoryService.BulkCreate(reqs)
	if err != nil {
		utils.HandleAppError(c, err, "failed to bulk create categories")
		return
	}
	utils.CreatedResponse(c, "categories created successfully", categories)
}

// BulkDelete removes multiple categories (admin only).
// @Summary      Bulk delete categories
// @Description  Deletes multiple categories by their IDs. Only accessible by users with the "admin" role.
// @Description  Cannot delete categories that have child categories – delete children first.
// @Tags         Admin Categories
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body object true "Category IDs to delete" SchemaExample({"ids":[1,2,3]})
// @Success      200 {object} utils.Response
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      403 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /categories/bulk [delete]
func (ctrl *CategoryController) BulkDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" validate:"required,min=1"`
	}
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}
	if err := ctrl.categoryService.BulkDelete(req.IDs); err != nil {
		utils.HandleAppError(c, err, "failed to bulk delete categories")
		return
	}
	utils.SuccessResponse(c, "categories deleted successfully", nil)
}
