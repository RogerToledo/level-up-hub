package activity

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/me/level-up-hub/backend/apperr"
	"github.com/me/level-up-hub/backend/config"
	"github.com/me/level-up-hub/backend/internal/pkg/identity"
	"github.com/me/level-up-hub/backend/internal/rest"
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
	activityID, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
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

	evidence, err := h.queries.AddEvidence(c.Request.Context(), activityID, userID, input.URL, input.Description)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, evidence, http.StatusCreated)
}

func (h *ActivityHandler) UpdateProgress(c *gin.Context) {
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

	var input struct {
		Progress int32 `json:"progress" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	if err := h.queries.UpdateProgress(c.Request.Context(), id, userID, input.Progress); err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ActivityHandler) Delete(c *gin.Context) {
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

	if err := h.queries.Delete(c.Request.Context(), id, userID); err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *ActivityHandler) GetDashboard(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	dashboard, err := h.queries.GetCareerDashboard(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, dashboard, http.StatusOK)
}

func (h *ActivityHandler) GetActivitiesEvidences(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	activities, err := h.queries.GetActivitiesEvidence(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, activities, http.StatusOK)
}

func (h *ActivityHandler) GetDetailedReport(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	report, err := h.queries.GetDetailedReport(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, report, http.StatusOK)
}

func (h *ActivityHandler) GetGapAnalysis(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	year := c.Query("year")
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrInvalidDate, err)
		return
	}

	gapAnalysis, err := h.queries.GetGapAnalysis(c.Request.Context(), userID, yearInt)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, gapAnalysis, http.StatusOK)
}

func (h *ActivityHandler) GetReadinessCheck(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	check, err := h.queries.GetCareerRadar(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, check, http.StatusOK)
}

func (h *ActivityHandler) GetCycleComparison(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	report, err := h.queries.GetCycleComparison(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, report, http.StatusOK)
}

func (h *ActivityHandler) DownloadReportPDF(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	activities, err := h.queries.GetDetailedReport(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	pdfBuffer, err := GenerateDossierPDF(activities)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=meu_dossie_carreira.pdf")
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", pdfBuffer.Len()))

	c.Writer.Write(pdfBuffer.Bytes())
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
