package account

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/me/level-up-hub/apperr"
	"github.com/me/level-up-hub/config"
	"github.com/me/level-up-hub/internal/rest"
)

type Handler struct {
	service *Service
	config  *config.Config
}

func NewHandler(s *Service, cfg *config.Config) *Handler {
	return &Handler{service: s, config: cfg}
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	token, err := h.service.Login(c.Request.Context(), req, h.config.JWTSecret)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrInvalidCredentials, nil)
		return
	}
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err.Error())
		return
	}

	rest.Send(c.Writer, gin.H{"token": token}, http.StatusOK)
}

func (h *Handler) Register(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	err := h.service.CreateUser(c.Request.Context(), req)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err.Error())
		return
	}

	rest.Send(c.Writer, fmt.Sprintf(apperr.OkCreate, apperr.UserPT), http.StatusCreated)
}

func (h *Handler) Update(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	id := c.Param("id")
	if id == "" {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, nil)
		return
	}

	err := h.service.UpdateUser(c.Request.Context(), id, req)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		rest.Error(c.Writer, http.StatusNotFound, fmt.Sprintf(apperr.ErrIsNotFound, apperr.UserPT), nil)
		return
	}
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err.Error())
		return
	}

	rest.Send(c.Writer, fmt.Sprintf(apperr.OkUpdate, apperr.UserPT), http.StatusOK)
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, nil)
		return
	}

	err := h.service.DeleteUser(c.Request.Context(), id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		rest.Error(c.Writer, http.StatusNotFound, fmt.Sprintf(apperr.ErrIsNotFound, apperr.UserPT), nil)
		return
	}
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err.Error())
		return
	}

	rest.Send(c.Writer, fmt.Sprintf(apperr.OkDelete, apperr.UserPT), http.StatusOK)
}

func (h *Handler) FindByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, nil)
		return
	}

	user, err := h.service.FindUserByID(c.Request.Context(), id)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		rest.Error(c.Writer, http.StatusNotFound, fmt.Sprintf(apperr.ErrIsNotFound, apperr.UserPT), nil)
		return
	}
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err.Error())
		return
	}

	rest.Send(c.Writer, user, http.StatusOK)
}

func (h *Handler) FindAll(c *gin.Context) {
	users, err := h.service.FindAllUsers(c.Request.Context())
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err.Error())
		return
	}

	rest.Send(c.Writer, users, http.StatusOK)
}
