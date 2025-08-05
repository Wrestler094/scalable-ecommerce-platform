package http

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/httphelper"

	"github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/delivery/http/dto"
	"github.com/Wrestler094/scalable-ecommerce-platform/catalog-service/internal/usecase"
)

type ProductHandler struct {
	productUC usecase.ProductUseCase
	validator httphelper.Validator
}

func NewProductHandler(uc usecase.ProductUseCase, validator httphelper.Validator) *ProductHandler {
	return &ProductHandler{productUC: uc, validator: validator}
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	req, err := httphelper.DecodeJSON[dto.CreateProductRequest](r, w)
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errFields := h.validator.Validate(req); errFields != nil {
		httphelper.RespondValidationErrors(w, errFields)
		return
	}

	id, err := h.productUC.CreateProduct(r.Context(), req)
	if err != nil {
		httphelper.RespondError(w, http.StatusInternalServerError, "failed to create product")
		return
	}

	httphelper.RespondJSON(w, http.StatusCreated, dto.CreateProductResponse{ID: id})
}

func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	p, err := h.productUC.GetProductByID(r.Context(), id)
	if err != nil {
		httphelper.RespondError(w, http.StatusInternalServerError, "failed to get product")
		return
	}

	httphelper.RespondJSON(w, http.StatusOK, dto.GetProductByIDResponse(*p))
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	req, err := httphelper.DecodeJSON[dto.UpdateProductRequest](r, w)
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errFields := h.validator.Validate(req); errFields != nil {
		httphelper.RespondValidationErrors(w, errFields)
		return
	}

	updated, err := h.productUC.UpdateProduct(r.Context(), id, req)
	if err != nil {
		httphelper.RespondError(w, http.StatusInternalServerError, "failed to update product")
		return
	}

	httphelper.RespondJSON(w, http.StatusOK, dto.UpdateProductResponse(updated))
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "invalid product id")
		return
	}

	if err := h.productUC.DeleteProduct(r.Context(), id); err != nil {
		httphelper.RespondError(w, http.StatusInternalServerError, "failed to delete product")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
