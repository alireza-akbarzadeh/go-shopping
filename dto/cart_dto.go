package dto

type AddItemResponse struct {
	BaseResponse
	Data CartItemData `json:"data"`
}

type CartItemData struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}

type GetCartResponse struct {
	BaseResponse
	Data GetCartData `json:"data"`
}

type GetCartData struct {
	Cart  CartData `json:"cart"`
	Total float64  `json:"total"`
}

type CartData struct {
	ID    uint             `json:"id"`
	Items []CartItemDetail `json:"items"`
}

type CartItemDetail struct {
	ID        uint    `json:"id"`
	ProductID uint    `json:"product_id"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
	Total     float64 `json:"total"`
}
