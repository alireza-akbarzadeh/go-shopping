package controllers

import (
	"errors"
	"strconv"

	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/middleware"
	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/services"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"github.com/gin-gonic/gin"
)

type ShipmentController struct {
	shipmentService services.ShipmentServiceInterface
}

func NewShipmentController(shipmentService services.ShipmentServiceInterface) *ShipmentController {
	return &ShipmentController{shipmentService: shipmentService}
}

// CreateShipment creates a new shipment (admin only) and enqueues background processing.
// @Summary      Create a shipment
// @Description  Creates a shipment record and triggers a background job to process it (e.g., call carrier API).
// @Tags         Admin Shipments
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
	var req struct {
		OrderID        uint   `json:"order_id" validate:"required"`
		Carrier        string `json:"carrier" validate:"required"`
		TrackingNumber string `json:"tracking_number"`
		AddressLine1   string `json:"address_line1" validate:"required"`
		AddressLine2   string `json:"address_line2"`
		City           string `json:"city" validate:"required"`
		State          string `json:"state"`
		PostalCode     string `json:"postal_code" validate:"required"`
		Country        string `json:"country" validate:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationErrorResponse(c, err.Error())
		return
	}

	// Basic validation – can be extended
	if req.OrderID == 0 {
		utils.ErrorResponse(c, 400, "order_id is required")
		return
	}

	shipment, err := ctrl.shipmentService.CreateShipment(
		req.OrderID,
		req.Carrier,
		req.TrackingNumber,
		req.AddressLine1,
		req.AddressLine2,
		req.City,
		req.State,
		req.PostalCode,
		req.Country,
	)
	if err != nil {
		var appErr *utils.AppError
		if errors.As(err, &appErr) {
			if appErr.Code == 404 {
				utils.NotFoundResponse(c, appErr.Message)
				return
			}
			if appErr.Code == 400 {
				utils.ErrorResponse(c, 400, appErr.Message)
				return
			}
		}
		utils.InternalServerErrorResponse(c, err, "failed to create shipment")
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
		var appErr *utils.AppError
		if errors.As(err, &appErr) && appErr.Code == 404 {
			utils.NotFoundResponse(c, constants.ErrNotFound)
			return
		}
		utils.InternalServerErrorResponse(c, err, "failed to fetch shipment")
		return
	}

	// Authorisation: admin can see all, users only their own
	if role != "admin" && shipment.UserID != userID {
		utils.ForbiddenResponse(c, constants.ErrForbidden)
		return
	}

	utils.SuccessResponse(c, constants.MsgFetchSuccess, shipment)
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

	orderID, err := strconv.ParseUint(c.Query("order_id"), 10, 64)
	if err != nil {
		utils.ErrorResponse(c, 400, "invalid order_id")
		return
	}

	// Verify order ownership (reuse service or direct DB check)
	// For simplicity we use the shipment service to fetch, but we need order ownership check.
	// We'll rely on the service to filter by userID.
	// Alternatively, we can query order directly. We'll assume the service enforces ownership.
	shipments, err := ctrl.shipmentService.GetShipmentsByOrderID(uint(orderID))
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
