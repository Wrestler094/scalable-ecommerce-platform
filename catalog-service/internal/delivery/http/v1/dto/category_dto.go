package dto

type Category struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// ====== CreateCategory ======

type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required"`
}

type CreateCategoryResponse struct {
	ID int64 `json:"id"`
}

// ====== GetCategory ======

type GetAllCategoriesResponse []Category

// ====== GetProductsByCategoryID ======

type GetProductsByCategoryIDResponse []Product