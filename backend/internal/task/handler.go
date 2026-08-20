package task

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/me/level-up-hub/backend/apperr"
	"github.com/me/level-up-hub/backend/config"
	"github.com/me/level-up-hub/backend/internal/pkg/identity"
	"github.com/me/level-up-hub/backend/internal/repository"
	"github.com/me/level-up-hub/backend/internal/rest"
)

type TaskHandler struct {
	service *Service
	cfg     *config.Config
}

func NewHandler(s *Service, cfg *config.Config) *TaskHandler {
	return &TaskHandler{service: s, cfg: cfg}
}

func (h *TaskHandler) Create(c *gin.Context) {
	var dto CreateTaskDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	task, err := h.service.Create(c.Request.Context(), dto)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, task, http.StatusCreated)
}

func (h *TaskHandler) Update(c *gin.Context) {
	taskID, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	var dto UpdateTaskDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	task, err := h.service.Update(c.Request.Context(), taskID, dto)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, task, http.StatusOK)
}

func (h *TaskHandler) Delete(c *gin.Context) {
	taskID, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), taskID); err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) ListByInitiative(c *gin.Context) {
	initiativeID, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	tasks, err := h.service.ListByInitiative(c.Request.Context(), initiativeID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}
	if tasks == nil {
		tasks = []repository.ListTasksByInitiativeRow{}
	}

	rest.Send(c.Writer, tasks, http.StatusOK)
}

func (h *TaskHandler) GetPillars(c *gin.Context) {
	taskID, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	pillars, err := h.service.GetTaskPillars(c.Request.Context(), taskID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}
	if pillars == nil {
		pillars = []repository.Pillar{}
	}

	rest.Send(c.Writer, pillars, http.StatusOK)
}

func (h *TaskHandler) AddEvidence(c *gin.Context) {
	taskID, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	var input struct {
		URL         string `json:"url" binding:"required,url"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	evidence, err := h.service.AddEvidence(c.Request.Context(), taskID, input.URL, input.Description)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, evidence, http.StatusCreated)
}

func (h *TaskHandler) ListEvidences(c *gin.Context) {
	taskID, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	evidences, err := h.service.ListEvidences(c.Request.Context(), taskID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}
	if evidences == nil {
		evidences = []repository.TaskEvidence{}
	}

	rest.Send(c.Writer, evidences, http.StatusOK)
}

func (h *TaskHandler) DeleteEvidence(c *gin.Context) {
	evidenceID, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	if err := h.service.DeleteEvidence(c.Request.Context(), evidenceID); err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) UpdateEvidence(c *gin.Context) {
	evidenceID, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	var input struct {
		URL         string `json:"url" binding:"required,url"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	evidence, err := h.service.UpdateEvidence(c.Request.Context(), evidenceID, input.URL, input.Description)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, evidence, http.StatusOK)
}
