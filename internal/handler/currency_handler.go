package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Tarifsiz/go-currency-api/internal/model"
	"github.com/Tarifsiz/go-currency-api/internal/service"
	"github.com/gin-gonic/gin"
)

// CurrencyHandler handles HTTP requests for currency operations
type CurrencyHandler struct {
	currencyService service.CurrencyServiceInterface
}

// NewCurrencyHandler creates a new currency handler instance
func NewCurrencyHandler(currencyService service.CurrencyServiceInterface) *CurrencyHandler {
	return &CurrencyHandler{
		currencyService: currencyService,
	}
}

// APIResponse represents the standard API response format
type APIResponse struct {
	Success   bool        `json:"success"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Message   string      `json:"message,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// PaginationResponse represents paginated API response
type PaginationResponse struct {
	Success    bool        `json:"success"`
	Data       interface{} `json:"data,omitempty"`
	Error      string      `json:"error,omitempty"`
	Message    string      `json:"message,omitempty"`
	Timestamp  time.Time   `json:"timestamp"`
	Pagination struct {
		Page    int   `json:"page"`
		Limit   int   `json:"limit"`
		Offset  int   `json:"offset"`
		Total   int64 `json:"total,omitempty"`
	} `json:"pagination,omitempty"`
}

// CreateCurrencyRequest represents the request body for creating a currency
type CreateCurrencyRequest struct {
	Code                string `json:"code" binding:"required,len=3"`
	Description         string `json:"description" binding:"required,max=255"`
	AmountDisplayFormat string `json:"amount_display_format,omitempty"`
	HtmlEncodedSymbol   string `json:"html_encoded_symbol,omitempty"`
	Factor              int    `json:"factor,omitempty"`
}

// UpdateCurrencyRequest represents the request body for updating a currency
type UpdateCurrencyRequest struct {
	Description         string `json:"description,omitempty"`
	AmountDisplayFormat string `json:"amount_display_format,omitempty"`
	HtmlEncodedSymbol   string `json:"html_encoded_symbol,omitempty"`
	Factor              int    `json:"factor,omitempty"`
}

// GetCurrencies handles GET /api/v1/currencies
func (h *CurrencyHandler) GetCurrencies(c *gin.Context) {
	// Parse query parameters
	page := h.getQueryInt(c, "page", 1)
	limit := h.getQueryInt(c, "limit", 50)
	search := c.Query("search")
	factor := h.getQueryInt(c, "factor", 0)
	
	// Calculate offset
	offset := (page - 1) * limit
	
	// Validate pagination parameters
	if limit > 100 {
		limit = 100 // Max limit
	}
	if limit < 1 {
		limit = 10 // Default limit
	}
	
	var currencies []*model.Currency
	var err error
	
	// Handle different query types
	if search != "" {
		currencies, err = h.currencyService.SearchCurrencies(c.Request.Context(), search)
	} else if factor > 0 {
		currencies, err = h.currencyService.GetCurrenciesByFactor(c.Request.Context(), factor)
	} else {
		currencies, err = h.currencyService.GetAllCurrencies(c.Request.Context(), limit, offset)
	}
	
	if err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to retrieve currencies", err)
		return
	}
	
	// Get total count for pagination (only for regular list, not search results)
	var total int64
	if search == "" && factor == 0 {
		total, _ = h.currencyService.GetCurrencyCount(c.Request.Context())
	}
	
	response := PaginationResponse{
		Success:   true,
		Data:      currencies,
		Timestamp: time.Now().UTC(),
	}
	
	response.Pagination.Page = page
	response.Pagination.Limit = limit
	response.Pagination.Offset = offset
	response.Pagination.Total = total
	
	c.JSON(http.StatusOK, response)
}

// GetCurrencyByCode handles GET /api/v1/currencies/:code
func (h *CurrencyHandler) GetCurrencyByCode(c *gin.Context) {
	code := strings.ToUpper(c.Param("code"))
	
	// Validate currency code format
	if len(code) != 3 {
		h.errorResponse(c, http.StatusBadRequest, "Invalid currency code format", nil)
		return
	}
	
	currency, err := h.currencyService.GetCurrencyByCode(c.Request.Context(), code)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Currency not found", err)
		return
	}
	
	h.successResponse(c, currency, "Currency retrieved successfully")
}

// CreateCurrency handles POST /api/v1/currencies
func (h *CurrencyHandler) CreateCurrency(c *gin.Context) {
	var req CreateCurrencyRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	// Convert to uppercase
	req.Code = strings.ToUpper(req.Code)
	
	// Create currency model
	currency := &model.Currency{
		Code:                req.Code,
		Description:         req.Description,
		AmountDisplayFormat: req.AmountDisplayFormat,
		HtmlEncodedSymbol:   req.HtmlEncodedSymbol,
		Factor:              req.Factor,
	}
	
	if err := h.currencyService.CreateCurrency(c.Request.Context(), currency); err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			h.errorResponse(c, http.StatusConflict, "Currency code already exists", err)
			return
		}
		h.errorResponse(c, http.StatusInternalServerError, "Failed to create currency", err)
		return
	}
	
	h.successResponse(c, currency, "Currency created successfully")
}

// UpdateCurrency handles PUT /api/v1/currencies/:code
func (h *CurrencyHandler) UpdateCurrency(c *gin.Context) {
	code := strings.ToUpper(c.Param("code"))
	
	// Validate currency code format
	if len(code) != 3 {
		h.errorResponse(c, http.StatusBadRequest, "Invalid currency code format", nil)
		return
	}
	
	var req UpdateCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.errorResponse(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}
	
	// Get existing currency
	currency, err := h.currencyService.GetCurrencyByCode(c.Request.Context(), code)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Currency not found", err)
		return
	}
	
	// Update fields if provided
	if req.Description != "" {
		currency.Description = req.Description
	}
	if req.AmountDisplayFormat != "" {
		currency.AmountDisplayFormat = req.AmountDisplayFormat
	}
	if req.HtmlEncodedSymbol != "" {
		currency.HtmlEncodedSymbol = req.HtmlEncodedSymbol
	}
	if req.Factor > 0 {
		currency.Factor = req.Factor
	}
	
	if err := h.currencyService.UpdateCurrency(c.Request.Context(), currency); err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to update currency", err)
		return
	}
	
	h.successResponse(c, currency, "Currency updated successfully")
}

// DeleteCurrency handles DELETE /api/v1/currencies/:code
func (h *CurrencyHandler) DeleteCurrency(c *gin.Context) {
	code := strings.ToUpper(c.Param("code"))
	
	// Validate currency code format
	if len(code) != 3 {
		h.errorResponse(c, http.StatusBadRequest, "Invalid currency code format", nil)
		return
	}
	
	// Get currency to get its ID
	currency, err := h.currencyService.GetCurrencyByCode(c.Request.Context(), code)
	if err != nil {
		h.errorResponse(c, http.StatusNotFound, "Currency not found", err)
		return
	}
	
	if err := h.currencyService.DeleteCurrency(c.Request.Context(), currency.ID); err != nil {
		h.errorResponse(c, http.StatusInternalServerError, "Failed to delete currency", err)
		return
	}
	
	h.successResponse(c, nil, "Currency deleted successfully")
}

// Helper methods

func (h *CurrencyHandler) getQueryInt(c *gin.Context, param string, defaultValue int) int {
	valueStr := c.Query(param)
	if valueStr == "" {
		return defaultValue
	}
	
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	
	return value
}

func (h *CurrencyHandler) successResponse(c *gin.Context, data interface{}, message string) {
	response := APIResponse{
		Success:   true,
		Data:      data,
		Message:   message,
		Timestamp: time.Now().UTC(),
	}
	
	statusCode := http.StatusOK
	if message == "Currency created successfully" {
		statusCode = http.StatusCreated
	}
	
	c.JSON(statusCode, response)
}

func (h *CurrencyHandler) errorResponse(c *gin.Context, statusCode int, message string, err error) {
	response := APIResponse{
		Success:   false,
		Error:     message,
		Timestamp: time.Now().UTC(),
	}
	
	// Log the actual error for debugging
	if err != nil {
		// In production, you'd want to use a proper logger
		println("Error:", err.Error())
	}
	
	c.JSON(statusCode, response)
}