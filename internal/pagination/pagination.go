package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// PaginationParams contains pagination parameters
type PaginationParams struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"page_size" json:"page_size"`
	Offset   int `json:"-"`
}

// PaginatedResponse contains paginated data and metadata
type PaginatedResponse struct {
	Data       interface{}        `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
}

// PaginationMetadata contains pagination information
type PaginationMetadata struct {
	CurrentPage  int   `json:"current_page"`
	PageSize     int   `json:"page_size"`
	TotalPages   int   `json:"total_pages"`
	TotalRecords int64 `json:"total_records"`
	HasNext      bool  `json:"has_next"`
	HasPrevious  bool  `json:"has_previous"`
}

// Pagination defaults and limits.
const (
	DefaultPage     = 1
	DefaultPageSize = 20
	MaxPageSize     = 100
)

// GetPaginationParams extracts and validates pagination parameters
func GetPaginationParams(c *gin.Context) PaginationParams {
	page, _ := strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(DefaultPage)))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(DefaultPageSize)))

	// Validations
	if page < 1 {
		page = DefaultPage
	}

	if pageSize < 1 {
		pageSize = DefaultPageSize
	}

	if pageSize > MaxPageSize {
		pageSize = MaxPageSize
	}

	offset := (page - 1) * pageSize

	return PaginationParams{
		Page:     page,
		PageSize: pageSize,
		Offset:   offset,
	}
}

// NewPaginatedResponse cria uma resposta paginada
func NewPaginatedResponse(data interface{}, params PaginationParams, totalRecords int64) PaginatedResponse {
	totalPages := int((totalRecords + int64(params.PageSize) - 1) / int64(params.PageSize))

	if totalPages == 0 {
		totalPages = 1
	}

	return PaginatedResponse{
		Data: data,
		Pagination: PaginationMetadata{
			CurrentPage:  params.Page,
			PageSize:     params.PageSize,
			TotalPages:   totalPages,
			TotalRecords: totalRecords,
			HasNext:      params.Page < totalPages,
			HasPrevious:  params.Page > 1,
		},
	}
}
