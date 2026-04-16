package activity

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/me/level-up-hub/backend/apperr"
	"github.com/me/level-up-hub/backend/config"
	"github.com/me/level-up-hub/backend/internal/email"
	"github.com/me/level-up-hub/backend/internal/pkg/identity"
	"github.com/me/level-up-hub/backend/internal/repository"
	"github.com/me/level-up-hub/backend/internal/rest"
)

type ActivityHandler struct {
	queries      *Service
	cfg          *config.Config
	emailService *email.Service
}

func NewHandler(s *Service, cfg *config.Config, emailService *email.Service) *ActivityHandler {
	return &ActivityHandler{
		queries:      s,
		cfg:          cfg,
		emailService: emailService,
	}
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

func (h *ActivityHandler) GetActivityEvidences(c *gin.Context) {
	activityID, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	evidences, err := h.queries.GetActivityEvidences(c.Request.Context(), activityID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	// Garante que sempre retorna um array, mesmo que vazio
	if evidences == nil {
		evidences = []repository.ActivityEvidence{}
	}

	rest.Send(c.Writer, evidences, http.StatusOK)
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

func (h *ActivityHandler) Update(c *gin.Context) {
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

	var dto UpdateActivityDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	if err := h.queries.Update(c.Request.Context(), id, userID, dto); err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, gin.H{"message": "Activity updated successfully"}, http.StatusOK)
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

func (h *ActivityHandler) List(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	activities, err := h.queries.ListActivities(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	// Garante que sempre retorna um array, mesmo que vazio
	if activities == nil {
		activities = []repository.ListUserActivitiesRow{}
	}

	rest.Send(c.Writer, activities, http.StatusOK)
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
	// Garante que sempre retorna um array, mesmo que vazio
	if activities == nil {
		activities = []repository.ListUserActivitiesWithEvidencesRow{}
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

	// Garante que sempre retorna um array, mesmo que vazio
	if report == nil {
		report = []repository.FindDetailedActivityReportRow{}
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

	// Garante que sempre retorna um array, mesmo que vazio
	if gapAnalysis == nil {
		gapAnalysis = []GapAnalysisResponse{}
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

	// Buscar dados completos do relatório (atividades + informações do usuário)
	reportData, err := h.queries.GetDetailedReportData(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	// Gerar PDF com dados completos
	pdfBuffer, err := GenerateDetailedDossierPDF(reportData)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=meu_dossie_carreira.pdf")
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Length", fmt.Sprintf("%d", pdfBuffer.Len()))

	c.Writer.Write(pdfBuffer.Bytes())
}

func (h *ActivityHandler) SendReportToManager(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	// Fetch user information including manager
	user, err := h.queries.repo.FindUserByID(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	// Check if user has a registered manager
	if !user.ManagerEmail.Valid || user.ManagerEmail.String == "" {
		rest.Error(c.Writer, http.StatusBadRequest, "Engineering manager not registered",
			"Please register your manager's email in your profile settings before sending the report.")
		return
	}

	// Fetch complete report data
	reportData, err := h.queries.GetDetailedReportData(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	// Gerar PDF
	pdfBuffer, err := GenerateDetailedDossierPDF(reportData)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	// Send email to manager
	managerName := "Manager"
	if user.ManagerName.Valid && user.ManagerName.String != "" {
		managerName = user.ManagerName.String
	}

	err = h.emailService.SendReportToManager(
		managerName,
		user.ManagerEmail.String,
		user.Username,
		user.Email,
		pdfBuffer.Bytes(),
	)

	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, "Error sending email", err)
		return
	}

	rest.Send(c.Writer, map[string]string{
		"message": "Report successfully sent to " + user.ManagerEmail.String,
		"status":  "success",
	}, http.StatusOK)
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
