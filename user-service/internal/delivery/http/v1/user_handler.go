package v1

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/httphelper"
	"github.com/Wrestler094/scalable-ecommerce-platform/pkg/logger"
	"github.com/Wrestler094/scalable-ecommerce-platform/user-service/internal/delivery/http/v1/dto"
	"github.com/Wrestler094/scalable-ecommerce-platform/user-service/internal/domain"
)

type UserHandler struct {
	userUC    domain.UserUseCase
	validator httphelper.Validator
	logger    logger.Logger
}

func NewUserHandler(userUC domain.UserUseCase, validator httphelper.Validator, logger logger.Logger) *UserHandler {
	return &UserHandler{
		userUC:    userUC,
		validator: validator,
		logger:    logger,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.Register"

	req, err := httphelper.DecodeJSON[dto.RegisterRequest](r, w)
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errFields := h.validator.Validate(req); errFields != nil {
		httphelper.RespondValidationErrors(w, errFields)
		return
	}

	userID, err := h.userUC.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			httphelper.RespondError(w, http.StatusConflict, "user already exists")
			return
		}

		h.logger.WithOp(op).
			WithRequestID(middleware.GetReqID(r.Context())).
			WithError(err).
			Error("failed to register user")

		httphelper.RespondError(w, http.StatusInternalServerError, "failed to register user")
		return
	}

	httphelper.RespondJSON(w, http.StatusOK, dto.RegisterResponse{UserID: userID})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.Login"

	req, err := httphelper.DecodeJSON[dto.LoginRequest](r, w)
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if errFields := h.validator.Validate(req); errFields != nil {
		httphelper.RespondValidationErrors(w, errFields)
		return
	}

	accessToken, refreshToken, err := h.userUC.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			httphelper.RespondError(w, http.StatusUnauthorized, "invalid email or password")
			return
		}

		h.logger.WithOp(op).
			WithRequestID(middleware.GetReqID(r.Context())).
			WithError(err).
			Error("failed to login user")

		httphelper.RespondError(w, http.StatusInternalServerError, "failed to login user")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	httphelper.RespondJSON(w, http.StatusOK, dto.LoginResponse{AccessToken: accessToken})
}

func (h *UserHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.Refresh"

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		httphelper.RespondError(w, http.StatusUnauthorized, "no refresh token")
		return
	}

	accessToken, err := h.userUC.Refresh(r.Context(), cookie.Value)
	if err != nil {
		h.logger.WithOp(op).
			WithRequestID(middleware.GetReqID(r.Context())).
			WithError(err).
			Warn("failed to refresh token")

		httphelper.RespondError(w, http.StatusUnauthorized, "invalid or expired refresh token")
		return
	}

	httphelper.RespondJSON(w, http.StatusOK, dto.RefreshResponse{AccessToken: accessToken})
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	const op = "UserHandler.Logout"

	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "no refresh token")
		return
	}

	if err := h.userUC.Logout(r.Context(), cookie.Value); err != nil {
		h.logger.WithOp(op).
			WithRequestID(middleware.GetReqID(r.Context())).
			WithError(err).
			Error("failed to logout user")

		httphelper.RespondError(w, http.StatusInternalServerError, "failed to logout user")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   -1,
		SameSite: http.SameSiteStrictMode,
	})

	httphelper.RespondJSON(w, http.StatusOK, dto.LogoutResponse{Success: true})
}
