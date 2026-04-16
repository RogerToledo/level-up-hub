package ladder

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/me/level-up-hub/backend/apperr"
	"github.com/me/level-up-hub/backend/config"
	"github.com/me/level-up-hub/backend/internal/repository"
	"github.com/me/level-up-hub/backend/internal/rest"
)

// LadderHandler handles HTTP requests for career ladder operations.
type LadderHandler struct {
	queries *Service
	cfg     *config.Config
}

// NewHandler creates a new ladder handler.
func NewHandler(s *Service, cfg *config.Config) *LadderHandler {
	return &LadderHandler{queries: s, cfg: cfg}
}

// Create handles the creation of a new ladder level.
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

// List handles retrieving all career ladder levels.
func (h *LadderHandler) List(c *gin.Context) {
	ladders, err := h.queries.ListAllLadders(c.Request.Context())
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}
	// Garante que sempre retorna um array, mesmo que vazio
	if ladders == nil {
		ladders = []repository.CareerLadder{}
	}
	rest.Send(c.Writer, ladders, http.StatusOK)
}
