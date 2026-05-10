package services

import (
	"errors"

	"github.com/alireza-akbarzadeh/shopping-platform/dto"
	"github.com/alireza-akbarzadeh/shopping-platform/models"
	"github.com/alireza-akbarzadeh/shopping-platform/utils"
	"gorm.io/gorm"
)

type addressService struct {
	db *gorm.DB
}

type AddressServiceInterface interface {
	Create(userID uint, req dto.CreateAddressRequest) (*models.Address, error)
	GetByID(id, userID uint) (*models.Address, error)
	Update(id, userID uint, req dto.UpdateAddressRequest) (*models.Address, error)
	Delete(id, userID uint) error
	List(userID uint) ([]models.Address, error)
	SetDefault(id, userID uint) error
	GetDefaultAddress(userID uint, addressType string) (*models.Address, error)
}

func NewAddressService(db *gorm.DB) AddressServiceInterface {
	return &addressService{db: db}
}

// Create a new address; if is_default is true, clear other default for that user+type.
func (s *addressService) Create(userID uint, req dto.CreateAddressRequest) (*models.Address, error) {
	// If this is the first address for the user, force is_default = true
	var count int64
	s.db.Model(&models.Address{}).Where("user_id = ?", userID).Count(&count)
	if count == 0 {
		req.IsDefault = true
	}

	// If setting as default, unset existing default for the same address_type (if 'both', unset both shipping and billing defaults? We'll handle simply)
	if req.IsDefault {
		if err := s.unsetDefault(userID, req.AddressType); err != nil {
			return nil, err
		}
	}

	address := &models.Address{
		UserID:        userID,
		AddressType:   req.AddressType,
		IsDefault:     req.IsDefault,
		RecipientName: req.RecipientName,
		Phone:         req.Phone,
		AddressLine1:  req.AddressLine1,
		AddressLine2:  req.AddressLine2,
		City:          req.City,
		State:         req.State,
		PostalCode:    req.PostalCode,
		Country:       req.Country,
		Instructions:  req.Instructions,
	}

	if err := s.db.Create(address).Error; err != nil {
		return nil, utils.ErrInternal(err)
	}
	return address, nil
}

// Helper: unset default for all addresses of a given type for the user
func (s *addressService) unsetDefault(userID uint, addressType string) error {
	return s.db.Model(&models.Address{}).
		Where("user_id = ? AND address_type IN (?, ?)", userID, addressType, "both").
		Where("address_type = ? OR address_type = 'both'", addressType).
		Update("is_default", false).Error
}

func (s *addressService) GetByID(id, userID uint) (*models.Address, error) {
	var addr models.Address
	err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&addr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, utils.ErrNotFound("address not found")
		}
		return nil, utils.ErrInternal(err)
	}
	return &addr, nil
}

func (s *addressService) Update(id, userID uint, req dto.UpdateAddressRequest) (*models.Address, error) {
	addr, err := s.GetByID(id, userID)
	if err != nil {
		return nil, err
	}

	if req.AddressType != nil {
		addr.AddressType = *req.AddressType
	}
	if req.IsDefault != nil && *req.IsDefault {
		// Unset previous default for the type
		if err := s.unsetDefault(userID, addr.AddressType); err != nil {
			return nil, err
		}
		addr.IsDefault = true
	} else if req.IsDefault != nil && !*req.IsDefault {
		addr.IsDefault = false
	}
	if req.RecipientName != nil {
		addr.RecipientName = *req.RecipientName
	}
	if req.Phone != nil {
		addr.Phone = *req.Phone
	}
	if req.AddressLine1 != nil {
		addr.AddressLine1 = *req.AddressLine1
	}
	if req.AddressLine2 != nil {
		addr.AddressLine2 = *req.AddressLine2
	}
	if req.City != nil {
		addr.City = *req.City
	}
	if req.State != nil {
		addr.State = *req.State
	}
	if req.PostalCode != nil {
		addr.PostalCode = *req.PostalCode
	}
	if req.Country != nil {
		addr.Country = *req.Country
	}
	if req.Instructions != nil {
		addr.Instructions = *req.Instructions
	}

	if err := s.db.Save(addr).Error; err != nil {
		return nil, utils.ErrInternal(err)
	}
	return addr, nil
}

func (s *addressService) Delete(id, userID uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.Address{})
	if result.Error != nil {
		return utils.ErrInternal(result.Error)
	}
	if result.RowsAffected == 0 {
		return utils.ErrNotFound("address not found")
	}
	return nil
}

func (s *addressService) List(userID uint) ([]models.Address, error) {
	var addresses []models.Address
	err := s.db.Where("user_id = ?", userID).Order("is_default DESC, created_at DESC").Find(&addresses).Error
	if err != nil {
		return nil, utils.ErrInternal(err)
	}
	return addresses, nil
}

func (s *addressService) SetDefault(id, userID uint) error {
	// First get the address to know its type
	addr, err := s.GetByID(id, userID)
	if err != nil {
		return err
	}
	// Unset existing default for that type
	if err := s.unsetDefault(userID, addr.AddressType); err != nil {
		return err
	}
	// Set this as default
	return s.db.Model(&models.Address{}).Where("id = ?", id).Update("is_default", true).Error
}

// GetDefaultAddress returns the default address of a given type for the user.
// If addressType is empty, it returns any default address (preferring 'both' or the first found).
// Returns nil, nil if no default address exists.
func (s *addressService) GetDefaultAddress(userID uint, addressType string) (*models.Address, error) {
	var addr models.Address
	query := s.db.Where("user_id = ? AND is_default = ?", userID, true)

	if addressType != "" {
		query = query.Where("address_type IN (?, ?)", addressType, "both")
	}

	err := query.First(&addr).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // no default address, not an error
		}
		return nil, utils.ErrInternal(err)
	}
	return &addr, nil
}
