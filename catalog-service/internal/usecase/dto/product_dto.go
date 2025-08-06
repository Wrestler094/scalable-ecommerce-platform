package dto

// CreateProductInput represents input for creating a product
type CreateProductInput struct {
	Name        string
	Description string
	Price       float64
	CategoryID  int64
}

// CreateProductOutput represents output for creating a product
type CreateProductOutput struct {
	ID int64
}

// UpdateProductInput represents input for updating a product
type UpdateProductInput struct {
	Name        *string
	Description *string
	Price       *float64
	CategoryID  *int64
}

// UpdateProductOutput represents output for updating a product
type UpdateProductOutput struct {
	ID          int64
	Name        string
	Description string
	Price       float64
	CategoryID  int64
}
