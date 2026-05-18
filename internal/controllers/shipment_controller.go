package controllers

import (
	"strconv"

	"github.com/alireza-akbarzadeh/luxe/internal/constants"
	"github.com/alireza-akbarzadeh/luxe/internal/middleware"
	"github.com/alireza-akbarzadeh/luxe/internal/models"
	"github.com/alireza-akbarzadeh/luxe/internal/services"
	"github.com/alireza-akbarzadeh/luxe/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ShipmentController struct {
	shipmentService services.ShipmentServiceInterface
	validate        *validator.Validate
}

func NewShipmentController(shipmentService services.ShipmentServiceInterface) *ShipmentController {
	return &ShipmentController{
		shipmentService: shipmentService,
		validate:        validator.New(),
	}
}

// CreateShipment creates a new shipment (admin only) and enqueues background processing.
// @Summary      Create a shipment
// @Description  Creates a shipment record and triggers a background job to process it (e.g., call carrier API).
// @Tags         Shipments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
//
//	@Param        request body object true "Shipment data" SchemaExample({
//	  "order_id":1,
//	  "carrier":"FedEx",
//	  "tracking_number":"123456789",
//	  "address_line1":"123 Main St",
//	  "city":"Springfield",
//	  "postal_code":"12345",
//	  "country":"USA"
//	})
//
// @Success      201 {object} utils.Response{data=models.Shipment}
// @Failure      400 {object} utils.Response
// @Failure      401 {object} utils.Response
// @Failure      403 {object} utils.Response
// @Failure      404 {object} utils.Response
// @Failure      500 {object} utils.Response
// @Router       /admin/shipments [post]
func (ctrl *ShipmentController) CreateShipment(c *gin.Context) {
	var req services.CreateShipmentRequest
	if !utils.BindAndValidate(c, &req, ctrl.validate) {
		return
	}

	shipment, err := ctrl.shipmentService.CreateShipment(req)
	if err != nil {
		utils.HandleAppError(c, err, "failed to create shipment")
		return
	}

	utils.CreatedResponse(c, constants.MsgCreateSuccess, shipment)
}

// GetShipment retrieves a shipment by ID (user sees own, admin sees any).
// @Summary      Get shipment by ID
// @Description  Returns a shipment. Admin can see any; users see only shipments belonging to their orders.
// @Tags         Shipments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Shipment ID"
// @Success      200  {object}  utils.Response{data=models.Shipment}
// @Failure      401  {object}  utils.Response
// @Failure      403  {object}  utils.Response
// @Failure      404  {object}  utils.Response
// @Router       /shipments/{id} [get]
func (ctrl *ShipmentController) GetShipment(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, constants.ErrUnauthorized)
		return
	}
	role, ok := middleware.GetUserRole(c)
	if !ok {
		utils.UnauthorizedResponse(c, constants.ErrUnauthorized)
		return
	}

	shipmentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid shipment id")
		return
	}

	shipment, err := ctrl.shipmentService.GetShipmentByID(uint(shipmentID))
	if err != nil {
		utils.HandleAppError(c, err, "failed to fetch shipment")
		return
	}

	// Authorisation: admin can see all, users only their own
	if role != "admin" && shipment.UserID != userID {
		utils.ForbiddenResponse(c, constants.ErrForbidden)
		return
	}

	utils.SuccessResponse(c, constants.MsgFetchSuccess, shipment)
}

type GetShipmentsByOrderRequest struct {
	OrderID uint `form:"order_id" validate:"required,gt=0"`
}

// GetShipmentsByOrder lists all shipments for a given order (user must own the order).
// @Summary      Get shipments for an order
// @Description  Returns all shipments belonging to an order (user must own the order).
// @Tags         Shipments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order_id  query  int  true  "Order ID"
// @Success      200       {object} utils.Response{data=[]models.Shipment}
// @Failure      400       {object} utils.Response
// @Failure      401       {object} utils.Response
// @Failure      403       {object} utils.Response
// @Failure      404       {object} utils.Response
// @Router       /shipments [get]
func (ctrl *ShipmentController) GetShipmentsByOrder(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		utils.UnauthorizedResponse(c, constants.ErrUnauthorized)
		return
	}

	var req GetShipmentsByOrderRequest
	if !utils.BindAndValidateQuery(c, &req, ctrl.validate) {
		return
	}

	// Verify order ownership (reuse service or direct DB check)
	// For simplicity we use the shipment service to fetch, but we need order ownership check.
	// We'll rely on the service to filter by userID.
	// Alternatively, we can query order directly. We'll assume the service enforces ownership.
	shipments, err := ctrl.shipmentService.GetShipmentsByOrderID(req.OrderID)
	if err != nil {
		utils.InternalServerErrorResponse(c, err, "failed to fetch shipments")
		return
	}

	// Filter by userID manually (or the service should do it)
	var userShipments []models.Shipment
	for _, s := range shipments {
		if s.UserID == userID {
			userShipments = append(userShipments, s)
		}
	}

	utils.SuccessResponse(c, constants.MsgFetchSuccess, userShipments)
}

// UpdateShipmentStatus updates a shipment's status (admin only).
// @Summary      Update shipment status (admin)
// @Description  Updates the status of a shipment and sends real-time notifications.
// @Tags         Shipments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id      path    int     true  "Shipment ID"
// @Param        request body    object  true  "Status update request"
// @Success      200     {object} utils.Response
// @Failure      400     {object} utils.Response
// @Failure      401     {object} utils.Response
// @Failure      403     {object} utils.Response
// @Failure      404     {object} utils.Response
// @Failure      500     {object} utils.Response
// @Router       /admin/shipments/{id}/status [put]
func (ctrl *ShipmentController) UpdateShipmentStatus(c *gin.Context) {
	shipmentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid shipment id")
		return
	}

	var req struct {
		Status string `json:"status" validate:"required,oneof=pending processing shipped delivered"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, err.Error())
		return
	}

	if err := ctrl.shipmentService.UpdateShipmentStatus(uint(shipmentID), req.Status); err != nil {
		utils.HandleAppError(c, err, "failed to update shipment status")
		return
	}

	utils.SuccessResponse(c, "shipment status updated successfully", nil)
}

// GetShippingProvider godoc
// @Summary      Get active shipping methods
// @Description  Returns all active shipping methods (public)
// @Tags         Shipping
// @Accept       json
// @Produce      json
// @Success      200 {object} utils.Response{data=[]models.ShippingMethod}
// @Router       /shipping-methods [get]
func (ctrl *ShipmentController) GetShippingProvider(c *gin.Context) {
	methods, err := ctrl.shipmentService.GetShippingProvider()
	if err != nil {
		utils.HandleAppError(c, err, "failed to fetch shipping methods")
		return
	}
	utils.SuccessResponse(c, "shipping methods retrieved", methods)
}
