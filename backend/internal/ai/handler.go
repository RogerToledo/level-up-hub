package ai

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/me/level-up-hub/backend/apperr"
	"github.com/me/level-up-hub/backend/internal/rest"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) Classify(c *gin.Context) {
	var req ClassifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	result, err := h.service.ClassifyTask(c.Request.Context(), req)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, "Erro ao classificar atividade", err.Error())
		return
	}

	rest.Send(c.Writer, result, http.StatusOK)
}
