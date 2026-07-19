package leveltarget

import (
	"fmt"
	"net/http"

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

func (h *Handler) Create(c *gin.Context) {
	var input repository.CreateLevelTargetParams
	if err := c.ShouldBindJSON(&input); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	lt, err := h.service.Create(c.Request.Context(), input)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, lt, http.StatusCreated)
}

func (h *Handler) FindByID(c *gin.Context) {
	id, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	lt, err := h.service.FindByID(c.Request.Context(), id)
	if err != nil {
		rest.Error(c.Writer, http.StatusNotFound, fmt.Sprintf(apperr.ErrIsNotFound, "meta de nivel"), nil)
		return
	}

	rest.Send(c.Writer, lt, http.StatusOK)
}

func (h *Handler) List(c *gin.Context) {
	targets, err := h.service.List(c.Request.Context())
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}
	if targets == nil {
		targets = []repository.ListLevelTargetsRow{}
	}

	rest.Send(c.Writer, targets, http.StatusOK)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	var input repository.UpdateLevelTargetParams
	if err := c.ShouldBindJSON(&input); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}
	input.ID = id

	lt, err := h.service.Update(c.Request.Context(), input)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, lt, http.StatusOK)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, fmt.Sprintf(apperr.OkDelete, "meta de nivel"), http.StatusOK)
}
