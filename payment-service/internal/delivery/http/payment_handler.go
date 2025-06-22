package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"payment-service/internal/delivery/http/dto"
	"pkg/authenticator"
	"pkg/httphelper"
	"pkg/logger"

	"payment-service/internal/domain"
)

type PaymentHandler struct {
	paymentUC domain.PaymentUseCase
	validator httphelper.Validator
	logger    logger.Logger
}

func NewPaymentHandler(paymentUC domain.PaymentUseCase, validator httphelper.Validator, logger logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentUC: paymentUC,
		validator: validator,
		logger:    logger,
	}
}

func (h *PaymentHandler) Pay(w http.ResponseWriter, r *http.Request) {
	const op = "paymentHandler.Pay"

	ctx := r.Context()
	userID, ok := authenticator.UserID(ctx)
	if !ok {
		httphelper.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	req, err := httphelper.DecodeJSON[dto.PayRequest](r, w)
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errFields := h.validator.Validate(req); errFields != nil {
		httphelper.RespondValidationErrors(w, errFields)
		return
	}

	payCommand := domain.PayCommand{
		UserID:         userID,
		OrderID:        req.OrderID,
		Amount:         req.Amount,
		IdempotencyKey: req.IdempotencyKey,
	}

	err = h.paymentUC.ProcessPayment(r.Context(), payCommand)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicatePayment) {
			http.Error(w, "duplicate payment", http.StatusConflict)
			return
		}

		if !errors.Is(err, domain.ErrIdempotencyRegistrationFailed) {
			h.logger.WithOp(op).
				WithRequestID(middleware.GetReqID(r.Context())).
				WithError(err).
				Error("payment failed", "command", payCommand)

			httphelper.RespondError(w, http.StatusInternalServerError, "failed to process payment")
			return
		}

		h.logger.WithOp(op).
			WithRequestID(middleware.GetReqID(r.Context())).
			WithError(err).
			Warn("idempotency registration failed", "command", payCommand)
	}

	httphelper.RespondJSON(w, http.StatusAccepted, dto.PayResponse{
		// TODO: Подумать над сообщением
		Message: "Payment accepted. Order status will be updated shortly.",
	})
}
