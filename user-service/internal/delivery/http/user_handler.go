package http

import (
	"net/http"

	"user-service/internal/delivery/http/dto"
	"user-service/internal/usecase"

	"pkg/httphelper"
)

type UserHandler struct {
	userUC    usecase.UserUseCase
	validator httphelper.Validator
}

func NewUserHandler(userUC usecase.UserUseCase, validator httphelper.Validator) *UserHandler {
	return &UserHandler{
		userUC:    userUC,
		validator: validator,
	}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
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
		httphelper.RespondError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	httphelper.RespondJSON(w, http.StatusOK, dto.RegisterResponse{UserID: userID})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
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
		httphelper.RespondError(w, http.StatusInternalServerError, "something went wrong")
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
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		httphelper.RespondError(w, http.StatusUnauthorized, "no refresh token")
		return
	}

	accessToken, err := h.userUC.Refresh(r.Context(), cookie.Value)
	if err != nil {
		httphelper.RespondError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	httphelper.RespondJSON(w, http.StatusOK, dto.RefreshResponse{AccessToken: accessToken})

}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		httphelper.RespondError(w, http.StatusBadRequest, "no refresh token")
		return
	}

	if err := h.userUC.Logout(r.Context(), cookie.Value); err != nil {
		httphelper.RespondError(w, http.StatusInternalServerError, "something went wrong")
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
