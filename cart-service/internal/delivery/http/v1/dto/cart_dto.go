package dto

type CartItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// ====== GetCart ======

type GetCartResponse struct {
	Items []CartItem `json:"items"`
}

// ====== AddItem ======

type AddItemRequest struct {
	ProductID int64 `json:"product_id" validate:"required,gt=0"`
}

type AddItemResponse CartItem

// ====== UpdateItem ======

type UpdateItemRequest struct {
	ProductID int64 `json:"product_id" validate:"required,gt=0"`
	Quantity  int   `json:"quantity" validate:"gte=0"`
}

type UpdateItemResponse CartItem

// ====== RemoveItem ======

type RemoveItemRequest struct {
	ProductID int64 `json:"product_id" validate:"required,gt=0"`
}

type RemoveItemResponse CartItem

// ====== ClearCart ======

type ClearCartResponse struct {
	Items []CartItem `json:"items"`
}