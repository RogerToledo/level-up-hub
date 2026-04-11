package activity

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/me/level-up-hub/apperr"
	"github.com/me/level-up-hub/config"
	"github.com/me/level-up-hub/internal/rest"
)

type ActivityHandler struct {
	queries *Service
	cfg     *config.Config
}

func NewHandler(s *Service, cfg *config.Config) *ActivityHandler {
	return &ActivityHandler{queries: s, cfg: cfg}
}

func (h *ActivityHandler) Create(c *gin.Context) {
	var dto CreateActivityDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		var details interface{}

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, fieldError := range validationErrors {
				errorMessages[fieldError.Field()] = getErrorMessage(fieldError)
			}
			details = errorMessages
		} else {
			details = err.Error()
		}

		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, details)
		return
	}

	err := h.queries.CreateCompleteActivity(c.Request.Context(), dto)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, fmt.Sprintf(apperr.OkCreateF, apperr.ActivityPT), http.StatusCreated)

}

func (h *ActivityHandler) AddEvidence(c *gin.Context) {
	activityID, _ := uuid.Parse(c.Param("id"))
	userID, _ := c.Get("user_id")

	var input struct {
		URL         string `json:"url" binding:"required,url"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "URL válida é obrigatória"})
		return
	}

	evidence, err := h.queries.AddEvidence(c.Request.Context(), activityID, userID.(uuid.UUID), input.URL, input.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao salvar evidência"})
		return
	}

	c.JSON(http.StatusCreated, evidence)
}

func (h *ActivityHandler) UpdateProgress(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	userID, _ := c.Get("user_id")

	var input struct {
		Progress int32 `json:"progress" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.queries.UpdateProgress(c.Request.Context(), id, userID.(uuid.UUID), input.Progress); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao atualizar"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ActivityHandler) Delete(c *gin.Context) {
	id, _ := uuid.Parse(c.Param("id"))
	userID, _ := c.Get("user_id")

	if err := h.queries.Delete(c.Request.Context(), id, userID.(uuid.UUID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao deletar"})
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ActivityHandler) GetDashboard(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in token"})
		return
	}

	userID := userIDValue.(uuid.UUID)

	dashboard, err := h.queries.GetCareerDashboard(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate dashboard"})
		return
	}

	c.JSON(http.StatusOK, dashboard)
}

func (h *ActivityHandler) GetActivitiesEvidences(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in token"})
		return
	}

	userID := userIDValue.(uuid.UUID)
	activities, err := h.queries.GetActivitiesEvidence(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch activities with evidences"})
		return
	}

	c.JSON(http.StatusOK, activities)
}

func (h *ActivityHandler) GetDetailedReport(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in token"})
		return
	}

	userID := userIDValue.(uuid.UUID)
	report, err := h.queries.GetDetailedReport(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch detailed report"})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *ActivityHandler) GetGapAnalysis(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in token"})
		return
	}

	userID := userIDValue.(uuid.UUID)
	year := c.Query("year")
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid year parameter"})
		return
	}

	gapAnalysis, err := h.queries.GetGapAnalysis(c.Request.Context(), userID, yearInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch gap analysis"})
		return
	}

	c.JSON(http.StatusOK, gapAnalysis)
}

func (h *ActivityHandler) GetReadinessCheck(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in token"})
		return
	}

	userID := userIDValue.(uuid.UUID)

	check, err := h.queries.GetCareerRadar(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to perform readiness check"})
		return
	}

	c.JSON(http.StatusOK, check)
}	

func (h *ActivityHandler) GetCycleComparison(c *gin.Context) {
    userIDValue, exists := c.Get("user_id")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in token"})
        return
    }

    userID := userIDValue.(uuid.UUID)

    report, err := h.queries.GetCycleComparison(c.Request.Context(), userID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, report)
}

func getErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("O campo '%s' é obrigatório", fe.Field())
	case "min":
		return fmt.Sprintf("O campo '%s' deve ser no mínimo %s", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("O campo '%s' deve ser no máximo %s", fe.Field(), fe.Param())
	case "email":
		return fmt.Sprintf("O campo '%s' deve ser um email válido", fe.Field())
	case "uuid":
		return fmt.Sprintf("O campo '%s' deve ser um UUID válido", fe.Field())
	default:
		return fmt.Sprintf("O campo '%s' é inválido", fe.Field())
	}
}
