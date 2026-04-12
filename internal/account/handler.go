package account

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/me/level-up-hub/apperr"
	"github.com/me/level-up-hub/config"
	"github.com/me/level-up-hub/internal/pagination"
	"github.com/me/level-up-hub/internal/pkg/identity"
	"github.com/me/level-up-hub/internal/rest"
)

// Handler handles HTTP requests for account operations.
type Handler struct {
	service *Service
	config  *config.Config
}

// NewHandler creates a new account handler.
func NewHandler(s *Service, cfg *config.Config) *Handler {
	return &Handler{service: s, config: cfg}
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user and return JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials  body      LoginRequest  true  "Login credentials"
// @Success      200          {object}  map[string]interface{}  "Login successful"
// @Failure      400          {object}  map[string]interface{}  "Invalid credentials"
// @Failure      401          {object}  map[string]interface{}  "Unauthorized"
// @Router       /login [post]
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

// Register godoc
// @Summary      Register new user
// @Description  Create a new user account
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserRequest  true  "User data"
// @Success      201   {object}  map[string]interface{}  "User created successfully"
// @Failure      400   {object}  map[string]interface{}  "Invalid data"
// @Failure      500   {object}  map[string]interface{}  "Internal error"
// @Router       /register [post]
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

// Update godoc
// @Summary      Update user
// @Description  Update user account information
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user  body      CreateUserRequest  true  "User data"
// @Success      200   {object}  map[string]interface{}  "User updated successfully"
// @Failure      400   {object}  map[string]interface{}  "Invalid data"
// @Router       /users [put]
func (h *Handler) Update(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	id, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	err = h.service.UpdateUser(c.Request.Context(), id, req)
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
	id, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	err = h.service.DeleteUser(c.Request.Context(), id)
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

// FindByID godoc
// @Summary      Get user by ID
// @Description  Returns information for a specific user
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID)"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}  "User found"
// @Failure      400  {object}  map[string]interface{}  "Invalid ID"
// @Failure      404  {object}  map[string]interface{}  "User not found"
// @Failure      401  {object}  map[string]interface{}  "Unauthorized"
// @Router       /users/{id} [get]
func (h *Handler) FindByID(c *gin.Context) {
	id, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
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

// FindAll godoc
// @Summary      List all users
// @Description  Returns paginated list of users (admin only)
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        page       query     int  false  "Page number"  default(1)
// @Param        page_size  query     int  false  "Items per page"  default(20)  maximum(100)
// @Security     BearerAuth
// @Success      200        {object}  map[string]interface{}  "User list"
// @Failure      401        {object}  map[string]interface{}  "Unauthorized"
// @Failure      500        {object}  map[string]interface{}  "Internal error"
// @Router       /users [get]
func (h *Handler) FindAll(c *gin.Context) {
	// Extract pagination parameters
	params := pagination.GetPaginationParams(c)

	// Fetch paginated users
	users, err := h.service.FindAllUsersPaginated(c.Request.Context(), params)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err.Error())
		return
	}

	// Count total users
	total, err := h.service.CountUsers(c.Request.Context())
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err.Error())
		return
	}

	// Create paginated response
	response := pagination.NewPaginatedResponse(users, params, total)
	rest.Send(c.Writer, response, http.StatusOK)
}
