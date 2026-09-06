package plan

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/me/level-up-hub/backend/apperr"
	"github.com/me/level-up-hub/backend/internal/pkg/identity"
	"github.com/me/level-up-hub/backend/internal/repository"
	"github.com/me/level-up-hub/backend/internal/rest"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Create(c *gin.Context) {
	var dto CreatePlanDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	if err := h.service.Create(c.Request.Context(), userID, dto); err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, fmt.Sprintf(apperr.OkCreateF, "plano"), http.StatusCreated)
}

func (h *Handler) Update(c *gin.Context) {
	id, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	var dto UpdatePlanDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	if err := h.service.Update(c.Request.Context(), id, userID, dto); err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, fmt.Sprintf(apperr.OkUpdate, "plano"), http.StatusOK)
}

func (h *Handler) Delete(c *gin.Context) {
	id, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id, userID); err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *Handler) FindByID(c *gin.Context) {
	id, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	plan, err := h.service.FindByID(c.Request.Context(), id, userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusNotFound, apperr.ErrIsNotFound, err)
		return
	}

	rest.Send(c.Writer, plan, http.StatusOK)
}

func (h *Handler) List(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	plans, err := h.service.ListPlans(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}
	if plans == nil {
		plans = []repository.ListUserPlansRow{}
	}

	rest.Send(c.Writer, plans, http.StatusOK)
}

func (h *Handler) Reorder(c *gin.Context) {
	id, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	var body struct {
		Position int32 `json:"position"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	if err := h.service.Reorder(c.Request.Context(), id, userID, body.Position); err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, fmt.Sprintf(apperr.OkUpdate, "plano"), http.StatusOK)
}

func (h *Handler) MoveUp(c *gin.Context) {
	id, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	plan, err := h.service.FindByID(c.Request.Context(), id, userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusNotFound, apperr.ErrIsNotFound, err)
		return
	}

	if plan.Position == 0 {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, fmt.Errorf("ja esta no topo"))
		return
	}

	if err := h.service.Reorder(c.Request.Context(), id, userID, plan.Position-1); err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, fmt.Sprintf(apperr.OkUpdate, "plano"), http.StatusOK)
}

func (h *Handler) MoveDown(c *gin.Context) {
	id, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	plan, err := h.service.FindByID(c.Request.Context(), id, userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusNotFound, apperr.ErrIsNotFound, err)
		return
	}

	count, err := h.service.repo.CountUserPlans(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	if plan.Position >= int32(count)-1 {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, fmt.Errorf("ja esta no final"))
		return
	}

	if err := h.service.Reorder(c.Request.Context(), id, userID, plan.Position+1); err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, fmt.Sprintf(apperr.OkUpdate, "plano"), http.StatusOK)
}