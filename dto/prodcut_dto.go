package dto

type BulkStockUpdate struct {
	ProductID uint `json:"product_id" validate:"required,gt=0"`
	Stock     int  `json:"stock" validate:"gte=0"`
}

type CreateProductRequest struct {
	Name              string   `json:"name" validate:"required,min=3,max=255"`
	Description       string   `json:"description,omitempty"`
	Price             float64  `json:"price" validate:"required,gte=0"`
	CompareAtPrice    *float64 `json:"compare_at_price,omitempty" validate:"omitempty,gte=0"`
	Cost              *float64 `json:"cost,omitempty" validate:"omitempty,gte=0"`
	SKU               string   `json:"sku" validate:"required,min=3,max=50"`
	Barcode           string   `json:"barcode,omitempty"`
	Stock             int      `json:"stock" validate:"gte=0"`
	LowStockThreshold int      `json:"low_stock_threshold,omitempty"`
	Weight            *float64 `json:"weight,omitempty" validate:"omitempty,gte=0"`
	IsDigital         bool     `json:"is_digital"`
	CategoryID        *uint    `json:"category_id,omitempty"`
	Images            []string `json:"images,omitempty"`
	Status            string   `json:"status" validate:"oneof=draft active inactive archived"`
	MetaTitle         string   `json:"meta_title,omitempty"`
	MetaDescription   string   `json:"meta_description,omitempty"`
}
type UpdateProductRequest struct {
	Name              *string   `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Description       *string   `json:"description,omitempty"`
	Price             *float64  `json:"price,omitempty" validate:"omitempty,gte=0"`
	CompareAtPrice    *float64  `json:"compare_at_price,omitempty" validate:"omitempty,gte=0"`
	Cost              *float64  `json:"cost,omitempty" validate:"omitempty,gte=0"`
	SKU               *string   `json:"sku,omitempty" validate:"omitempty,min=3,max=50"`
	Barcode           *string   `json:"barcode,omitempty"`
	Stock             *int      `json:"stock,omitempty" validate:"omitempty,gte=0"`
	LowStockThreshold *int      `json:"low_stock_threshold,omitempty"`
	Weight            *float64  `json:"weight,omitempty" validate:"omitempty,gte=0"`
	IsDigital         *bool     `json:"is_digital,omitempty"`
	CategoryID        *uint     `json:"category_id,omitempty"`
	Images            *[]string `json:"images,omitempty"`
	Status            *string   `json:"status,omitempty" validate:"omitempty,oneof=draft active inactive archived"`
	MetaTitle         *string   `json:"meta_title,omitempty"`
	MetaDescription   *string   `json:"meta_description,omitempty"`
}

type BulkDeleteProductsRequest struct {
	ProductIDs []uint `json:"product_ids" validate:"required,min=1"`
}
