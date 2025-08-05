package v1

import (
	"net/http"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/authenticator"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/httphelper"

	"github.com/Wrestler094/scalable-ecommerce-platform/cart-service/internal/delivery/http/v1/dto"
	"github.com/Wrestler094/scalable-ecommerce-platform/cart-service/internal/domain"
)

type CartHandler struct {
	cartUC    domain.CartUseCase
	validator httphelper.Validator
}

func NewCartHandler(uc domain.CartUseCase, validator httphelper.Validator) *CartHandler {
	return &CartHandler{
		cartUC:    uc,
		validator: validator,
	}
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := authenticator.UserID(ctx)
	if !ok {
		httphelper.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	items, err := h.cartUC.GetCart(ctx, userID)
	if err != nil {
		httphelper.RespondError(w, http.StatusInternalServerError, "failed to get cart")
		return
	}

	respItems := make([]dto.CartItem, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, dto.CartItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}

	resp := dto.GetCartResponse{
		Items: respItems,
	}

	httphelper.RespondJSON(w, http.StatusOK, resp)
}

func (h *CartHandler) AddItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := authenticator.UserID(ctx)
	if !ok {
		httphelper.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	req, err := httphelper.DecodeJSON[dto.AddItemRequest](r, w)
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errFields := h.validator.Validate(req); errFields != nil {
		httphelper.RespondValidationErrors(w, errFields)
		return
	}

	if err := h.cartUC.AddItem(ctx, userID, req.ProductID, 1); err != nil {
		httphelper.RespondError(w, http.StatusInternalServerError, "failed to add item")
		return
	}

	resp := dto.AddItemResponse{
		ProductID: req.ProductID,
		Quantity:  1,
	}

	httphelper.RespondJSON(w, http.StatusOK, resp)
}

func (h *CartHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := authenticator.UserID(ctx)
	if !ok {
		httphelper.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	req, err := httphelper.DecodeJSON[dto.UpdateItemRequest](r, w)
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errFields := h.validator.Validate(req); errFields != nil {
		httphelper.RespondValidationErrors(w, errFields)
		return
	}

	if err := h.cartUC.UpdateItem(ctx, userID, req.ProductID, req.Quantity); err != nil {
		httphelper.RespondError(w, http.StatusInternalServerError, "failed to update item")
		return
	}

	resp := dto.UpdateItemResponse{
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	httphelper.RespondJSON(w, http.StatusOK, resp)
}

func (h *CartHandler) RemoveItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := authenticator.UserID(ctx)
	if !ok {
		httphelper.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	req, err := httphelper.DecodeJSON[dto.RemoveItemRequest](r, w)
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errFields := h.validator.Validate(req); errFields != nil {
		httphelper.RespondValidationErrors(w, errFields)
		return
	}

	if err := h.cartUC.RemoveItem(ctx, userID, req.ProductID); err != nil {
		httphelper.RespondError(w, http.StatusInternalServerError, "failed to remove item")
		return
	}

	resp := dto.RemoveItemResponse{
		ProductID: req.ProductID,
		Quantity:  0,
	}

	httphelper.RespondJSON(w, http.StatusOK, resp)
}

func (h *CartHandler) ClearCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, ok := authenticator.UserID(ctx)
	if !ok {
		httphelper.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	if err := h.cartUC.ClearCart(ctx, userID); err != nil {
		httphelper.RespondError(w, http.StatusInternalServerError, "failed to clear cart")
		return
	}

	resp := dto.ClearCartResponse{
		Items: []dto.CartItem{},
	}

	httphelper.RespondJSON(w, http.StatusOK, resp)
}