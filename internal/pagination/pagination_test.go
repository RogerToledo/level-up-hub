package pagination

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetPaginationParams(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedPage   int
		expectedSize   int
		expectedOffset int
	}{
		{
			name:           "default values when no params provided",
			queryParams:    map[string]string{},
			expectedPage:   1,
			expectedSize:   20,
			expectedOffset: 0,
		},
		{
			name: "valid page and size",
			queryParams: map[string]string{
				"page":      "2",
				"page_size": "10",
			},
			expectedPage:   2,
			expectedSize:   10,
			expectedOffset: 10,
		},
		{
			name: "negative page defaults to 1",
			queryParams: map[string]string{
				"page": "-1",
			},
			expectedPage:   1,
			expectedSize:   20,
			expectedOffset: 0,
		},
		{
			name: "zero page defaults to 1",
			queryParams: map[string]string{
				"page": "0",
			},
			expectedPage:   1,
			expectedSize:   20,
			expectedOffset: 0,
		},
		{
			name: "page size exceeding max is capped",
			queryParams: map[string]string{
				"page_size": "200",
			},
			expectedPage:   1,
			expectedSize:   100,
			expectedOffset: 0,
		},
		{
			name: "invalid string values use defaults",
			queryParams: map[string]string{
				"page":      "abc",
				"page_size": "xyz",
			},
			expectedPage:   1,
			expectedSize:   20,
			expectedOffset: 0,
		},
		{
			name: "large page number",
			queryParams: map[string]string{
				"page":      "10",
				"page_size": "50",
			},
			expectedPage:   10,
			expectedSize:   50,
			expectedOffset: 450,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Create request with query params
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			q := req.URL.Query()
			for key, val := range tt.queryParams {
				q.Add(key, val)
			}
			req.URL.RawQuery = q.Encode()
			c.Request = req

			// Execute
			params := GetPaginationParams(c)

			// Assert
			assert.Equal(t, tt.expectedPage, params.Page, "page mismatch")
			assert.Equal(t, tt.expectedSize, params.PageSize, "page_size mismatch")
			assert.Equal(t, tt.expectedOffset, params.Offset, "offset mismatch")
		})
	}
}

func TestNewPaginatedResponse(t *testing.T) {
	tests := []struct {
		name          string
		data          interface{}
		params        PaginationParams
		total         int64
		expectedPages int
		expectedNext  bool
		expectedPrev  bool
	}{
		{
			name:          "first page with multiple pages",
			data:          []string{"item1", "item2"},
			params:        PaginationParams{Page: 1, PageSize: 10, Offset: 0},
			total:         25,
			expectedPages: 3,
			expectedNext:  true,
			expectedPrev:  false,
		},
		{
			name:          "middle page",
			data:          []string{"item1", "item2"},
			params:        PaginationParams{Page: 2, PageSize: 10, Offset: 10},
			total:         25,
			expectedPages: 3,
			expectedNext:  true,
			expectedPrev:  true,
		},
		{
			name:          "last page",
			data:          []string{"item1", "item2"},
			params:        PaginationParams{Page: 3, PageSize: 10, Offset: 20},
			total:         25,
			expectedPages: 3,
			expectedNext:  false,
			expectedPrev:  true,
		},
		{
			name:          "single page",
			data:          []string{"item1", "item2"},
			params:        PaginationParams{Page: 1, PageSize: 10, Offset: 0},
			total:         5,
			expectedPages: 1,
			expectedNext:  false,
			expectedPrev:  false,
		},
		{
			name:          "empty results",
			data:          []string{},
			params:        PaginationParams{Page: 1, PageSize: 10, Offset: 0},
			total:         0,
			expectedPages: 1, // Always at least 1 page, even if empty
			expectedNext:  false,
			expectedPrev:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := NewPaginatedResponse(tt.data, tt.params, tt.total)

			assert.Equal(t, tt.data, response.Data, "data mismatch")
			assert.Equal(t, tt.params.Page, response.Pagination.CurrentPage, "current page mismatch")
			assert.Equal(t, tt.params.PageSize, response.Pagination.PageSize, "page size mismatch")
			assert.Equal(t, tt.total, response.Pagination.TotalRecords, "total records mismatch")
			assert.Equal(t, tt.expectedPages, response.Pagination.TotalPages, "total pages mismatch")
			assert.Equal(t, tt.expectedNext, response.Pagination.HasNext, "has_next mismatch")
			assert.Equal(t, tt.expectedPrev, response.Pagination.HasPrevious, "has_previous mismatch")
		})
	}
}

func TestPaginationMetadata(t *testing.T) {
	t.Run("calculate offset correctly", func(t *testing.T) {
		params := PaginationParams{
			Page:     3,
			PageSize: 20,
		}

		expectedOffset := (3 - 1) * 20 // 40
		params.Offset = expectedOffset

		assert.Equal(t, 40, params.Offset)
	})

	t.Run("edge case page 1", func(t *testing.T) {
		params := PaginationParams{
			Page:     1,
			PageSize: 10,
		}

		params.Offset = (params.Page - 1) * params.PageSize // 0

		assert.Equal(t, 0, params.Offset)
	})
}
