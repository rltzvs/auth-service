package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"auth/internal/controller"
	"auth/internal/entity"
)

type AuthHandler struct {
	service controller.AuthService
	logger  *zap.Logger
}

func NewAuthHandler(service controller.AuthService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AuthHandler) sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	}); err != nil {
		h.logger.Error("failed to encode error response", zap.Error(err))
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var user entity.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := user.Validate()
	if err != nil {
		h.logger.Error("failed to validate user", zap.Error(err))
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.Register(r.Context(), user)
	if err != nil {
		h.logger.Error("failed to register user", zap.Error(err))
		// Проверяем тип ошибки для правильного статус-кода
		if errors.Is(err, entity.ErrUserAlreadyExists) {
			h.sendError(w, http.StatusConflict, "Email already exists")
		} else {
			h.sendError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	h.logger.Info("user registered successfully", zap.String("email", user.Email))
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	}); err != nil {
		h.logger.Error("failed to encode response", zap.Error(err))
		return
	}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var user entity.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		h.logger.Error("failed to decode request body", zap.Error(err))
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	err := user.Validate()
	if err != nil {
		h.logger.Error("failed to validate user", zap.Error(err))
		h.sendError(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := h.service.Login(r.Context(), user)
	if err != nil {
		h.logger.Error("failed to login user", zap.Error(err))
		switch {
		case errors.Is(err, entity.ErrUserNotFound):
			h.sendError(w, http.StatusUnauthorized, "Invalid email or password")
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			h.sendError(w, http.StatusUnauthorized, "Invalid email or password")
		default:
			h.sendError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	h.logger.Info("user logged in successfully", zap.String("email", user.Email))
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	}); err != nil {
		h.logger.Error("failed to encode response", zap.Error(err))
		return
	}
}
