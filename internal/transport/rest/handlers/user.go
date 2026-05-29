package handlers

import (
	"docintel/internal/app"
	"docintel/internal/domain"
	"docintel/internal/transport/rest/response"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type UserHandler struct {
	AuthService domain.AuthService
}

func NewUserHandler(authService domain.AuthService) *UserHandler {
	return &UserHandler{
		AuthService: authService,
	}
}

func (h *UserHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req domain.RegisterUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	err = h.AuthService.Register(ctx, req.Email, req.Name, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value") {
			response.WriteError(w, http.StatusBadRequest, errors.New(domain.ErrUserAlreadyExists))
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resp := response.Response{
		Message: "User registered successfully",
	}

	response.WriteJSON(w, http.StatusCreated, resp)
}

func (h *UserHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req domain.LoginUserRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, errors.New("invalid request body"))
		return
	}

	user, err := h.AuthService.Login(ctx, req.Email, req.Password)
	if err != nil {
		if err.Error() == app.InvalidCredentialsErr {
			response.WriteError(w, http.StatusUnauthorized, errors.New(app.InvalidCredentialsErr))
			return
		}
		response.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	resp := response.Response{
		Message: "Login successful",
		Data: map[string]interface{}{
			"user": user,
		},
	}

	response.WriteJSON(w, http.StatusOK, resp)
}
