package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/alireza-akbarzadeh/shopping-platform/constants"
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
// @Param        request body services.CreateProductRequest true "Product details" SchemaExample({"name":"Laptop","price":999.99,"description":"High performance laptop","stock":10})
// @Success      201 {object} utils.Response{data=object{product=object{id=uint,name=string,price=float64,description=string,stock=int,created_at=string,updated_at=string}}}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      409 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /products [post]
func (ctrl *ProductController) Create(c *gin.Context) {
	var req services.CreateProductRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request payload: "+err.Error())
		return
	}

	if err := ctrl.validate.Struct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}

	product, err := ctrl.productService.Create(req)
	if err != nil {
		var appErr *utils.AppError
		if errors.As(err, &appErr) && appErr.Code == 409 {
			utils.ConflictResponse(c, appErr.Message)
			return
		}
		utils.InternalServerErrorResponse(c, err, "failed to create product")
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
// @Param        request body services.UpdateProductRequest true "Product update data" SchemaExample({"name":"Updated Laptop","price":1099.99,"description":"Even better performance","stock":5})
// @Success      200 {object} utils.Response{data=object{product=object{id=uint,name=string,price=float64,description=string,stock=int,updated_at=string}}}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Failure      409 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /products/{id} [put]
func (ctrl *ProductController) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	var req services.UpdateProductRequest
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid product id")
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, "Invalid request payload: "+err.Error())
		return
	}

	if err := ctrl.validate.Struct(req); err != nil {
		utils.HandleValidationError(c, err)
		return
	}
	product, err := ctrl.productService.Update(uint(id), req)
	if err != nil {
		var appErr *utils.AppError
		if errors.As(err, &appErr) {
			switch appErr.Code {
			case http.StatusNotFound:
				utils.NotFoundResponse(c, appErr.Message)
			case http.StatusConflict:
				utils.ConflictResponse(c, appErr.Message)
			default:
				utils.InternalServerErrorResponse(c, err, "failed to update product")
			}
		}
		utils.InternalServerErrorResponse(c, err, "failed to update product")
		return
	}
	utils.SuccessResponse(c, constants.MsgUpdateSuccess, product)
}
