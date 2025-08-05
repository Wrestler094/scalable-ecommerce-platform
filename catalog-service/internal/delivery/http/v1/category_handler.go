package v1

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/httphelper"

	"github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/delivery/http/v1/dto"
	"github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/usecase"
)

type CategoryHandler struct {
	categoryUC usecase.CategoryUseCase
	productUC  usecase.ProductUseCase
	validator  httphelper.Validator
}

func NewCategoryHandler(
	categoryUC usecase.CategoryUseCase,
	productUC usecase.ProductUseCase,
	validator httphelper.Validator,
) *CategoryHandler {
	return &CategoryHandler{
		categoryUC: categoryUC,
		productUC:  productUC,
		validator:  validator,
	}
}

func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	req, err := httphelper.DecodeJSON[dto.CreateCategoryRequest](r, w)
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errFields := h.validator.Validate(req); errFields != nil {
		httphelper.RespondValidationErrors(w, errFields)
		return
	}

	id, err := h.categoryUC.CreateCategory(r.Context(), req.Name)
	if err != nil {
		httphelper.RespondError(w, http.StatusInternalServerError, "failed to create category")
		return
	}

	httphelper.RespondJSON(w, http.StatusCreated, dto.CreateCategoryResponse{ID: id})
}

func (h *CategoryHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryUC.ListCategories(r.Context())
	if err != nil {
		httphelper.RespondError(w, http.StatusInternalServerError, "failed to get categories")
		return
	}

	var resp dto.GetAllCategoriesResponse
	for _, c := range categories {
		resp = append(resp, dto.Category{
			ID:   c.ID,
			Name: c.Name,
		})
	}

	httphelper.RespondJSON(w, http.StatusOK, resp)
}

func (h *CategoryHandler) GetProductsByCategoryID(w http.ResponseWriter, r *http.Request) {
	categoryID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "invalid category id")
		return
	}

	products, err := h.productUC.ListByCategory(r.Context(), categoryID)
	if err != nil {
		httphelper.RespondError(w, http.StatusInternalServerError, "failed to get products by category")
		return
	}

	var resp dto.GetProductsByCategoryIDResponse
	for _, p := range products {
		resp = append(resp, dto.Product{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			CategoryID:  p.CategoryID,
		})
	}

	httphelper.RespondJSON(w, http.StatusOK, resp)
}
