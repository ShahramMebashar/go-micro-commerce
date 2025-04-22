package api

import (
	"encoding/json"
	"microservice/services/product-service/internal/infrastructure/validator"
	"net/http"
)

// Response types
const (
	ResponseTypeSuccess = "success"
	ResponseTypeError   = "error"
)

// Standard HTTP status messages
var statusMessages = map[int]string{
	http.StatusOK:                  "ok",
	http.StatusCreated:             "Created",
	http.StatusAccepted:            "Accepted",
	http.StatusNoContent:           "No Content",
	http.StatusBadRequest:          "Bad Request",
	http.StatusUnauthorized:        "Unauthorized",
	http.StatusForbidden:           "Forbidden",
	http.StatusNotFound:            "Not Found",
	http.StatusConflict:            "Conflict",
	http.StatusInternalServerError: "Internal Server Error",
}

// APIResponse is the standard response format for all API endpoints
type APIResponse struct {
	Type    string      `json:"type"`              // "success" or "error"
	Message string      `json:"message,omitempty"` // Human-readable message
	Data    interface{} `json:"data,omitempty"`    // Response data for success
	Errors  interface{} `json:"errors,omitempty"`  // Error details for failures
}

// PaginatedResponse adds pagination metadata to the response
type PaginatedResponse struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

// Pagination holds pagination metadata
type Pagination struct {
	CurrentPage int  `json:"current_page"`
	PerPage     int  `json:"per_page"`
	TotalItems  int  `json:"total_items"`
	TotalPages  int  `json:"total_pages"`
	HasPrevPage bool `json:"has_prev_page"`
	HasNextPage bool `json:"has_next_page"`
}

// NewPagination creates a new Pagination struct
func NewPagination(currentPage, perPage, totalItems int) Pagination {
	totalPages := totalItems / perPage
	if totalItems%perPage > 0 {
		totalPages++
	}

	return Pagination{
		CurrentPage: currentPage,
		PerPage:     perPage,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		HasPrevPage: currentPage > 1,
		HasNextPage: currentPage < totalPages,
	}
}

// RespondWithJSON sends a JSON response with the given status code
func RespondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	// Get default message for status code or use "Success" as fallback
	message, ok := statusMessages[statusCode]
	if !ok {
		message = "Success"
	}

	// Determine response type based on status code
	responseType := ResponseTypeSuccess
	if statusCode >= 400 {
		responseType = ResponseTypeError
	}

	response := APIResponse{
		Type:    responseType,
		Message: message,
	}

	// Set data or errors based on response type
	if responseType == ResponseTypeSuccess {
		response.Data = data
	} else {
		response.Errors = data
	}

	// Marshal response to JSON
	jsonData, err := json.Marshal(response)
	if err != nil {
		// If marshaling fails, send a simple error response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"type":"error","message":"Failed to encode response"}`))
		return
	}

	// Set content type and status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)
}

// RespondWithError sends an error response
func RespondWithError(w http.ResponseWriter, message string, statusCode int) {
	RespondWithJSON(w, statusCode, map[string]string{"message": message})
}

// RespondWithValidationErrors sends validation errors
func RespondWithValidationErrors(w http.ResponseWriter, errors []validator.ValidationError) {
	RespondWithJSON(w, http.StatusBadRequest, errors)
}

// RespondWithPagination sends a paginated response
func RespondWithPagination(w http.ResponseWriter, items interface{}, page, perPage, totalItems int) {
	pagination := NewPagination(page, perPage, totalItems)
	paginatedResponse := PaginatedResponse{
		Items:      items,
		Pagination: pagination,
	}
	RespondWithJSON(w, http.StatusOK, paginatedResponse)
}
