package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/alireza-akbarzadeh/shopping-platform/constants"
	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/tasks"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"gorm.io/gorm"
)

type ShipmentServiceInterface interface {
	CreateShipment(orderID uint, carrier, trackingNumber string,
		addressLine1, addressLine2, city, state, postalCode, country string) (*models.Shipment, error)
	GetShipmentByID(id uint) (*models.Shipment, error)
	GetShipmentsByOrderID(orderID uint) ([]models.Shipment, error)
	UpdateShipmentStatus(id uint, status string) error
}

type shipmentService struct {
	db         *gorm.DB
	workerPool *tasks.WorkerPool
}

func NewShipmentService(db *gorm.DB, workerPool *tasks.WorkerPool) ShipmentServiceInterface {
	return &shipmentService{db: db, workerPool: workerPool}
}

// CreateShipment creates a shipment record and enqueues a background job.
func (s *shipmentService) CreateShipment(orderID uint, carrier, trackingNumber string,
	addressLine1, addressLine2, city, state, postalCode, country string) (*models.Shipment, error) {

	var order models.Order
	if err := s.db.First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound(constants.ErrOrderNotFound)
		}
		return nil, utils.ErrInternal(err)
	}

	shipment := &models.Shipment{
		OrderID:        orderID,
		UserID:         order.UserID,
		Carrier:        carrier,
		TrackingNumber: trackingNumber,
		Status:         "pending",
		AddressLine1:   addressLine1,
		AddressLine2:   addressLine2,
		City:           city,
		State:          state,
		PostalCode:     postalCode,
		Country:        country,
	}

	if err := s.db.Create(shipment).Error; err != nil {
		return nil, utils.ErrInternal(err)
	}

	// Enqueue background job (e.g., call carrier API, generate label)
	job := tasks.Job{
		ID:      fmt.Sprintf("shipment_%d", shipment.ID),
		Payload: shipment.ID,
		Handler: s.processShipment,
	}
	s.workerPool.Enqueue(job)

	return shipment, nil
}

// processShipment is the job handler called by the worker pool.
func (s *shipmentService) processShipment(payload interface{}) error {
	shipmentID, ok := payload.(uint)
	if !ok {
		return fmt.Errorf("invalid payload type")
	}

	// Simulate external carrier API call (label generation, tracking update)
	time.Sleep(2 * time.Second)

	// Update status to 'processing' (or 'shipped' after successful API call)
	if err := s.db.Model(&models.Shipment{}).Where("id = ?", shipmentID).
		Update("status", "processing").Error; err != nil {
		return err
	}

	utils.Log.Infof("Shipment %d processed in background", shipmentID)
	return nil
}

// GetShipmentByID retrieves a shipment by its ID.
func (s *shipmentService) GetShipmentByID(id uint) (*models.Shipment, error) {
	var shipment models.Shipment
	if err := s.db.First(&shipment, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound(constants.ErrShipmentNotFound)
		}
		return nil, utils.ErrInternal(err)
	}
	return &shipment, nil
}

// GetShipmentsByOrderID returns all shipments associated with an order.
func (s *shipmentService) GetShipmentsByOrderID(orderID uint) ([]models.Shipment, error) {
	var shipments []models.Shipment
	if err := s.db.Where("order_id = ?", orderID).Find(&shipments).Error; err != nil {
		return nil, utils.ErrInternal(err)
	}
	return shipments, nil
}

// UpdateShipmentStatus manually updates a shipment's status (admin only).
func (s *shipmentService) UpdateShipmentStatus(id uint, status string) error {
	result := s.db.Model(&models.Shipment{}).Where("id = ?", id).Update("status", status)
	if result.Error != nil {
		return utils.ErrInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.ErrNotFound(constants.ErrShipmentNotFound)
	}
	return nil
}
