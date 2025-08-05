package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/authenticator"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/httphelper"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"

	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/delivery/http/dto"
	"github.com/Wrestler094/scalable-ecommerce-platform/order-service/internal/domain"
)

type OrderHandler struct {
	orderUC   domain.OrderUseCase
	validator httphelper.Validator
	logger    logger.Logger
}

func NewOrderHandler(orderUC domain.OrderUseCase, validator httphelper.Validator, logger logger.Logger) *OrderHandler {
	return &OrderHandler{
		orderUC:   orderUC,
		validator: validator,
		logger:    logger,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	const op = "orderHandler.CreateOrder"

	ctx := r.Context()
	userID, ok := authenticator.UserID(ctx)
	if !ok {
		httphelper.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	req, err := httphelper.DecodeJSON[dto.CreateOrderRequest](r, w)
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errFields := h.validator.Validate(req); errFields != nil {
		httphelper.RespondValidationErrors(w, errFields)
		return
	}

	items := req.ToDomainItems()
	order, paymentURL, err := h.orderUC.CreateOrder(ctx, userID, items)
	if err != nil {
		h.logger.WithOp(op).WithError(err).Error("Failed to create order")
		httphelper.RespondError(w, http.StatusInternalServerError, "failed to create order")
		return
	}

	resp := dto.CreateOrderResponse{
		Order:      dto.FromOrder(order),
		PaymentURL: paymentURL,
	}
	httphelper.RespondJSON(w, http.StatusCreated, resp)
}

func (h *OrderHandler) GetOrdersList(w http.ResponseWriter, r *http.Request) {
	const op = "orderHandler.GetOrdersList"

	ctx := r.Context()
	userID, ok := authenticator.UserID(ctx)
	if !ok {
		httphelper.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orders, err := h.orderUC.ListOrdersByUser(ctx, userID)
	if err != nil {
		h.logger.WithOp(op).WithError(err).Error("Failed to get orders list")
		httphelper.RespondError(w, http.StatusInternalServerError, "failed to get orders list")
		return
	}

	resp := dto.GetOrdersListResponse{
		Orders: dto.FromOrders(orders),
	}

	httphelper.RespondJSON(w, http.StatusOK, resp)
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	const op = "orderHandler.GetOrderByID"

	ctx := r.Context()
	userID, ok := authenticator.UserID(ctx)
	if !ok {
		httphelper.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	orderID := chi.URLParam(r, "id")
	if _, err := uuid.Parse(orderID); err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "invalid uuid")
		return
	}

	order, err := h.orderUC.GetOrderByUUID(ctx, orderID)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			httphelper.RespondError(w, http.StatusNotFound, "order not found")
			return
		}

		h.logger.WithOp(op).WithError(err).Error("Failed to get order", "order_id", orderID)
		httphelper.RespondError(w, http.StatusInternalServerError, "failed to get order")
		return
	}

	if order.UserID != userID {
		httphelper.RespondError(w, http.StatusForbidden, "forbidden")
		return
	}

	resp := dto.GetOrderByIDResponse{
		Order: dto.FromOrder(order),
	}
	httphelper.RespondJSON(w, http.StatusOK, resp)
}
