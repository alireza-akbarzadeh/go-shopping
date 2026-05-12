package controllers

import (
	"strconv"

	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/dto"
	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/services"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ProductController struct {
	productService services.ProductServiceInterface
	validate       *validator.Validate
}

func NewProductController(productServices services.ProductServiceInterface) *ProductController {
	return &ProductController{
		productService: productServices,
		validate:       validator.New(),
	}
}

// Create creates a new product.
// @Summary      Create product
// @Description  Add a new product to the catalog
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.CreateProductRequest true "Product details"
// @Success      201 {object} utils.Response{data=models.Product}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      409 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /products [post]
func (ctrl *ProductController) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}

	product, err := ctrl.productService.Create(req)
	if err != nil {
		utils.HandleAppError(c, err, "failed to create product")
		return
	}

	utils.CreatedResponse(c, constants.MsgCreateSuccess, product)
}

// Update updates an existing product.
// @Summary      Update product
// @Description  Update an existing product by its ID
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Product ID"
// @Param        request body dto.UpdateProductRequest true "Product update data"
// @Success      200 {object} utils.Response{data=models.Product}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Failure      409 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /products/{id} [put]
func (ctrl *ProductController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	var req dto.UpdateProductRequest
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid product id")
		return
	}
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}
	product, err := ctrl.productService.Update(uint(id), req)
	if err != nil {
		utils.HandleAppError(c, err, "failed to update product")
		return
	}
	utils.SuccessResponse(c, constants.MsgUpdateSuccess, product)
}

// Delete removes a product by ID.
// @Summary      Delete product
// @Description  Delete an existing product permanently from the catalog
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Product ID"
// @Success      200 {object} utils.Response
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /products/{id} [delete]
func (ctrl *ProductController) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid product id")
		return
	}

	if err := ctrl.productService.Delete(uint(id)); err != nil {
		utils.HandleAppError(c, err, "failed to delete product")
		return
	}

	utils.SuccessResponse[any](c, constants.MsgDeleteSuccess, nil)
}

// GetOne retrieves a product by ID or slug.
// @Summary      Get product by ID or slug
// @Description  Fetch a single product using either its numeric ID or slug (URL-friendly name)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        identifier path string true "Product identifier (ID or slug)"
// @Success      200  {object}  dto.ProductSingleResponse
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /products/{identifier} [get]
func (ctrl *ProductController) GetOne(c *gin.Context) {
	identifier := c.Param("identifier")
	id, err := strconv.ParseUint(identifier, 10, 64)
	var product *models.Product
	if err == nil {
		product, err = ctrl.productService.GetByID(uint(id))
	} else {
		product, err = ctrl.productService.GetBySlug(identifier)
	}
	if err != nil {
		utils.HandleAppError(c, err, "failed to fetch product")
		return
	}
	utils.SuccessResponse(c, constants.MsgFetchSuccess, product)
}

type ProductListFilters struct {
	Status     string  `form:"status"`
	Name       string  `form:"name"`
	CategoryID uint    `form:"category_id"`
	MinPrice   float64 `form:"min_price"`
	MaxPrice   float64 `form:"max_price"`
	IsDigital  *bool   `form:"is_digital"`
}

// List retrieves a paginated list of products with optional filters.
// @Summary      List products
// @Description  Get products with pagination and filtering by status, name, category, price range, and digital flag.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Items per page (default 20, max 100)" default(20)
// @Param        offset query int false "Number of items to skip (default 0)" default(0)
// @Param        status query string false "Product status (e.g., active, draft, archived)" Example(active)
// @Param        name query string false "Filter by product name (partial match)" Example(laptop)
// @Param        category_id query int false "Filter by category ID" Example(5)
// @Param        min_price query number false "Minimum price filter" Example(10.99)
// @Param        max_price query number false "Maximum price filter" Example(999.99)
// @Param        is_digital query boolean false "Filter digital products (true/false)" Example(true)
// @Success      200  {object}  dto.ProductListResponse
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /products [get]
func (ctrl *ProductController) List(c *gin.Context) {
	limit := constants.DefaultLimit
	if l, err := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(constants.DefaultLimit))); err == nil && l >= constants.MinLimit {
		limit = l
	}
	if limit > constants.MaxLimit {
		limit = constants.MaxLimit
	}

	offset := constants.MinOffset
	if o, err := strconv.Atoi(c.DefaultQuery("offset", strconv.Itoa(constants.MinOffset))); err == nil && o >= constants.MinOffset {
		offset = o
	}

	// Bind filters
	var filters ProductListFilters
	if !utils.BindAndValidateQuery(c, &filters, ctrl.validate) {
		return
	}

	// Convert to map (if service expects map) or pass struct directly
	filterMap := make(map[string]interface{})
	if filters.Status != "" {
		filterMap["status"] = filters.Status
	}
	if filters.Name != "" {
		filterMap["name"] = filters.Name
	}
	if filters.CategoryID != 0 {
		filterMap["category_id"] = filters.CategoryID
	}
	if filters.MinPrice != 0 {
		filterMap["min_price"] = filters.MinPrice
	}
	if filters.MaxPrice != 0 {
		filterMap["max_price"] = filters.MaxPrice
	}
	if filters.IsDigital != nil {
		filterMap["is_digital"] = *filters.IsDigital
	}

	products, total, err := ctrl.productService.List(limit, offset, filterMap)
	if err != nil {
		utils.InternalServerErrorResponse(c, err, "failed to list products")
		return
	}

	data := gin.H{
		"products": products,
		"total":    total,
		"limit":    limit,
		"offset":   offset,
	}
	utils.SuccessResponse(c, constants.MsgFetchSuccess, data)
}

// BulkCreate creates multiple products at once (admin only).
// @Summary      Bulk create products
// @Description  Create multiple products in a single request (admin only)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body []dto.CreateProductRequest true "Array of products"
// @Success      201 {object} utils.Response{data=[]models.Product}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      403 {object} utils.Response
// @Router       /products/bulk [post]
func (ctrl *ProductController) BulkCreate(c *gin.Context) {
	var reqs []dto.CreateProductRequest
	if !utils.BindAndValidate(c, &reqs, ctrl.validate) {
		return
	}
	if len(reqs) == 0 {
		utils.ErrorResponse(c, 400, "no products provided")
		return
	}
	products, err := ctrl.productService.BulkCreate(reqs)
	if err != nil {
		utils.HandleAppError(c, err, "failed to bulk create products")
		return
	}
	utils.CreatedResponse(c, constants.MsgCreateSuccess, products)
}

// BulkDelete removes multiple products at once (admin only).
// @Summary      Bulk delete products
// @Description  Soft delete multiple products by their IDs (admin only)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.BulkDeleteProductsRequest true "Product IDs to delete"
// @Success      200 {object} utils.Response
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      403 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Router       /products/bulk [delete]
func (ctrl *ProductController) BulkDelete(c *gin.Context) {
	var req dto.BulkDeleteProductsRequest
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}
	if err := ctrl.productService.BulkDelete(req.ProductIDs); err != nil {
		utils.HandleAppError(c, err, "failed to delete products")
		return
	}
	utils.SuccessResponse[any](c, "products deleted successfully", nil)
}
