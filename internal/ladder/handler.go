package ladder

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/me/level-up-hub/apperr"
	"github.com/me/level-up-hub/config"
	"github.com/me/level-up-hub/internal/repository"
	"github.com/me/level-up-hub/internal/rest"
)

type LadderHandler struct {
	queries *Service
	cfg     *config.Config
}

func NewHandler(s *Service, cfg *config.Config) *LadderHandler {
	return &LadderHandler{queries: s, cfg: cfg}
}

func (h *LadderHandler) Create(c *gin.Context) {
	var input repository.CreateLadderLevelParams

	if err := c.ShouldBindJSON(&input); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	err := h.queries.CreateLadderLevel(c.Request.Context(), input)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, fmt.Sprintf(apperr.OkCreate, apperr.LadderPT), http.StatusCreated)
}
