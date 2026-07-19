package initiative

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

type InitiativeHandler struct {
	service      *Service
	cfg          *config.Config
	emailService *email.Service
}

func NewHandler(s *Service, cfg *config.Config, emailService *email.Service) *InitiativeHandler {
	return &InitiativeHandler{
		service:      s,
		cfg:          cfg,
		emailService: emailService,
	}
}

func (h *InitiativeHandler) Create(c *gin.Context) {
	var dto CreateInitiativeDTO

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

	err := h.service.CreateCompleteInitiative(c.Request.Context(), dto)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, fmt.Sprintf(apperr.OkCreateF, "iniciativa"), http.StatusCreated)
}

func (h *InitiativeHandler) Update(c *gin.Context) {
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

	var dto UpdateInitiativeDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	if err := h.service.Update(c.Request.Context(), id, userID, dto); err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, fmt.Sprintf(apperr.OkUpdate, "iniciativa"), http.StatusOK)
}

func (h *InitiativeHandler) Delete(c *gin.Context) {
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

func (h *InitiativeHandler) List(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	initiatives, err := h.service.ListInitiatives(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	if initiatives == nil {
		initiatives = []repository.ListUserInitiativesRow{}
	}

	rest.Send(c.Writer, initiatives, http.StatusOK)
}

func (h *InitiativeHandler) GetPillars(c *gin.Context) {
	initiativeID, err := identity.ValidateIDParam(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusBadRequest, apperr.ErrBadRequest, err)
		return
	}

	pillars, err := h.service.GetInitiativePillars(c.Request.Context(), initiativeID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	if pillars == nil {
		pillars = []repository.Pillar{}
	}

	rest.Send(c.Writer, pillars, http.StatusOK)
}

func (h *InitiativeHandler) GetDashboard(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	dashboard, err := h.service.GetCareerDashboard(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, dashboard, http.StatusOK)
}

func (h *InitiativeHandler) GetDetailedReport(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	report, err := h.service.GetDetailedReport(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	if report == nil {
		report = []repository.FindDetailedInitiativeReportRow{}
	}

	rest.Send(c.Writer, report, http.StatusOK)
}

func (h *InitiativeHandler) GetGapAnalysis(c *gin.Context) {
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

	gapAnalysis, err := h.service.GetGapAnalysis(c.Request.Context(), userID, yearInt)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	if gapAnalysis == nil {
		gapAnalysis = []GapAnalysisResponse{}
	}

	rest.Send(c.Writer, gapAnalysis, http.StatusOK)
}

func (h *InitiativeHandler) GetReadinessCheck(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	check, err := h.service.GetCareerRadar(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, check, http.StatusOK)
}

func (h *InitiativeHandler) GetCycleComparison(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	report, err := h.service.GetCycleComparison(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	rest.Send(c.Writer, report, http.StatusOK)
}

func (h *InitiativeHandler) DownloadReportPDF(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	reportData, err := h.service.GetDetailedReportData(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

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

func (h *InitiativeHandler) SendReportToManager(c *gin.Context) {
	userID, err := identity.GetUserIDFromContext(c)
	if err != nil {
		rest.Error(c.Writer, http.StatusUnauthorized, apperr.ErrUnauthorized, err)
		return
	}

	user, err := h.service.repo.FindUserByID(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	if !user.ManagerEmail.Valid || user.ManagerEmail.String == "" {
		rest.Error(c.Writer, http.StatusBadRequest, "Engineering manager not registered",
			"Please register your manager's email in your profile settings before sending the report.")
		return
	}

	reportData, err := h.service.GetDetailedReportData(c.Request.Context(), userID)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

	pdfBuffer, err := GenerateDetailedDossierPDF(reportData)
	if err != nil {
		rest.Error(c.Writer, http.StatusInternalServerError, apperr.ErrInternalServerError, err)
		return
	}

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
		return fmt.Sprintf("O campo '%s' e obrigatorio", fe.Field())
	case "min":
		return fmt.Sprintf("O campo '%s' deve ser no minimo %s", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("O campo '%s' deve ser no maximo %s", fe.Field(), fe.Param())
	case "email":
		return fmt.Sprintf("O campo '%s' deve ser um email valido", fe.Field())
	case "uuid":
		return fmt.Sprintf("O campo '%s' deve ser um UUID valido", fe.Field())
	default:
		return fmt.Sprintf("O campo '%s' e invalido", fe.Field())
	}
}
