package controllers

import (
	"net/http"
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
// @Success      201 {object} dto.ProductSingleResponse
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

	c.JSON(http.StatusCreated, dto.ProductSingleResponse{
		BaseResponse: dto.BaseResponse{
			Success: true,
			Message: constants.MsgCreateSuccess,
			Code:    http.StatusCreated,
		},
		Data: dto.ProductSingleData{Product: product},
	})
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
// @Success      200 {object} dto.ProductSingleResponse
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
	c.JSON(http.StatusOK, dto.ProductSingleResponse{
		BaseResponse: dto.BaseResponse{
			Success: true,
			Message: constants.MsgUpdateSuccess,
			Code:    http.StatusOK,
		},
		Data: dto.ProductSingleData{Product: product},
	})
}

// Delete removes a product by ID.
// @Summary      Delete product
// @Description  Delete an existing product permanently from the catalog
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Product ID"
// @Success      200 {object} dto.EmptyResponse
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

	c.JSON(http.StatusOK, dto.EmptyResponse{
		BaseResponse: dto.BaseResponse{
			Success: true,
			Message: constants.MsgDeleteSuccess,
			Code:    http.StatusOK,
		},
	})
}

// GetOne retrieves a product by ID or slug.
// @Summary      Get product by ID or slug
// @Description  Fetch a single product using either its numeric ID or slug (URL-friendly name)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "Product identifier (ID or slug)"
// @Success      200  {object}  dto.ProductSingleResponse
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /products/{id} [get]
func (ctrl *ProductController) GetOne(c *gin.Context) {
	identifier := c.Param("id")
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
	c.JSON(http.StatusOK, dto.ProductSingleResponse{
		BaseResponse: dto.BaseResponse{
			Success: true,
			Message: constants.MsgFetchSuccess,
			Code:    http.StatusOK,
		},
		Data: dto.ProductSingleData{Product: product},
	})
}

// List retrieves a paginated list of products with optional filters.
// @Summary      List products
// @Description  Get products with pagination and filtering by status, name, category, price range, rating, reviews count, digital flag, new flag, and sorting.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit query int false "Items per page (default 20, max 100)" default(20)
// @Param        offset query int false "Number of items to skip (default 0)" default(0)
// @Param        status query string false "Product status (active, draft, archived)" Example(active)
// @Param        name query string false "Filter by exact product name" Example(laptop)
// @Param        sku query string false "Filter by SKU" Example(SKU-123)
// @Param        category_id query int false "Filter by category ID" Example(5)
// @Param        min_price query number false "Minimum price filter" Example(10.99)
// @Param        max_price query number false "Maximum price filter" Example(999.99)
// @Param        min_rating query number false "Minimum rating (0.0 to 5.0)" Example(4.0)
// @Param        max_rating query number false "Maximum rating" Example(5.0)
// @Param        min_reviews query int false "Minimum number of reviews" Example(10)
// @Param        max_reviews query int false "Maximum number of reviews" Example(1000)
// @Param        is_digital query bool false "Filter digital or physical products" Example(true)
// @Param        is_new query bool false "Filter newly added products (first 30 days)" Example(true)
// @Param        sort query string false "Sort order: rating_desc, rating_asc, newest, reviews_desc, price_asc, price_desc" Example(rating_desc)
// @Success      200  {object}  dto.ProductListResponse
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /products [get]
func (ctrl *ProductController) List(c *gin.Context) {
	// Parse limit & offset
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

	// Bind filters directly from query string
	var filters dto.ProductListFilters
	if !utils.BindAndValidateQuery(c, &filters, ctrl.validate) {
		return
	}

	// Optional: run validator if you have validation rules on the struct
	// if err := ctrl.validate.Struct(filters); err != nil { ... }

	// Call service with struct (no map conversion needed)
	products, total, err := ctrl.productService.List(limit, offset, filters)
	if err != nil {
		utils.InternalServerErrorResponse(c, err, "failed to list products")
		return
	}

	c.JSON(http.StatusOK, dto.ProductListResponse{
		BaseResponse: dto.BaseResponse{
			Success: true,
			Message: constants.MsgFetchSuccess,
			Code:    http.StatusOK,
		},
		Data: dto.ProductListData{
			Products: products,
			Total:    total,
			Limit:    limit,
			Offset:   offset,
		},
	})
}

// BulkCreate creates multiple products at once (admin only).
// @Summary      Bulk create products
// @Description  Create multiple products in a single request (admin only)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body []dto.CreateProductRequest true "Array of products"
// @Success      201 {object} dto.ProductListResponse
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
	c.JSON(http.StatusCreated, dto.ProductListResponse{
		BaseResponse: dto.BaseResponse{
			Success: true,
			Message: constants.MsgCreateSuccess,
			Code:    http.StatusCreated,
		},
		Data: dto.ProductListData{
			Products: products,
			Total:    int64(len(products)),
			Limit:    len(products),
			Offset:   0,
		},
	})
}

// BulkDelete removes multiple products at once (admin only).
// @Summary      Bulk delete products
// @Description  Soft delete multiple products by their IDs (admin only)
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body dto.BulkDeleteProductsRequest true "Product IDs to delete"
// @Success      200 {object} dto.EmptyResponse
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
	c.JSON(http.StatusOK, dto.EmptyResponse{
		BaseResponse: dto.BaseResponse{
			Success: true,
			Message: "products deleted successfully",
			Code:    http.StatusOK,
		},
	})
}

// GetRelated retrieves products related to a given product (same category, excluding itself).
// @Summary      Get related products
// @Description  Fetch a list of products from the same category as the specified product, ordered by rating and review count.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Product ID"
// @Param        limit query int false "Number of related products (default 4, max 10)"
// @Success      200 {object} dto.ProductListResponse
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /products/{id}/related [get]
func (ctrl *ProductController) GetRelated(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid product id")
		return
	}

	limit := 4
	if l, err := strconv.Atoi(c.DefaultQuery("limit", "4")); err == nil && l > 0 {
		if l > 10 {
			limit = 10
		} else {
			limit = l
		}
	}

	products, err := ctrl.productService.GetRelated(uint(id), limit)
	if err != nil {
		utils.HandleAppError(c, err, "failed to fetch related products")
		return
	}

	c.JSON(http.StatusOK, dto.ProductListResponse{
		BaseResponse: dto.BaseResponse{
			Success: true,
			Message: constants.MsgFetchSuccess,
			Code:    http.StatusOK,
		},
		Data: dto.ProductListData{
			Products: products,
			Total:    int64(len(products)),
			Limit:    limit,
			Offset:   0,
		},
	})
}
