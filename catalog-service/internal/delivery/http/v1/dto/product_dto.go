package dto

type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  int64   `json:"category_id"`
}

// ====== CreateProduct ======

type CreateProductRequest struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	CategoryID  int64   `json:"category_id" validate:"required,gt=0"`
}

type CreateProductResponse struct {
	ID int64 `json:"id"`
}

// ====== GetProductByID ======

type GetProductByIDResponse Product

// ====== UpdateProduct ======

type UpdateProductRequest struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty,gt=0"`
	CategoryID  *int64   `json:"category_id,omitempty,gt=0"`
}

type UpdateProductResponse Product
