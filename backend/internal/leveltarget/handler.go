package leveltarget

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/me/level-up-hub/backend/apperr"
	"github.com/me/level-up-hub/backend/config"
	"github.com/me/level-up-hub/backend/internal/pkg/identity"
	"github.com/me/level-up-hub/backend/internal/repository"
	"github.com/me/level-up-hub/backend/internal/rest"
)

type Handler struct {
	service *Service
	cfg     *config.Config
}

func NewHandler(s *Service, cfg *config.Config) *Handler {
	return &Handler{service: s, cfg: cfg}
}

// SetTarget sets or updates the user's target level for a given year
func (h *Handler) SetTarget(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	var input struct {
		Target repository.LadderLevel `json:"target" binding:"required"`
		Year   int32                  `json:"year"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	// Default to current year if not provided
	if input.Year == 0 {
		input.Year = int32(time.Now().Year())
	}

	lt, err := h.service.Upsert(c.Request.Context(), userID, input.Target, input.Year)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, lt, http.StatusOK)
}

// GetTarget returns the user's target for a given year (defaults to current year)
func (h *Handler) GetTarget(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	yearStr := c.Query("year")
	year := int32(time.Now().Year())
	if yearStr != "" {
		y, err := strconv.Atoi(yearStr)
		if err != nil {
			rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
			return
		}
		year = int32(y)
	}

	lt, err := h.service.FindByUserAndYear(c.Request.Context(), userID, year)
	if err != nil {
		// No target set yet — return empty
		rest.Send(c.Writer, nil, http.StatusOK)
		return
	}

	rest.Send(c.Writer, lt, http.StatusOK)
}

// ListTargets returns all targets for the authenticated user
func (h *Handler) ListTargets(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	targets, err := h.service.ListByUser(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}
	if targets == nil {
		targets = []repository.ListLevelTargetsByUserRow{}
	}

	rest.Send(c.Writer, targets, http.StatusOK)
}

// DeleteTarget removes a target for the authenticated user
func (h *Handler) DeleteTarget(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	id, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}
