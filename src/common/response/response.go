package response

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
)

// StandardResponse represents the standard API response structure
type StandardResponse struct {
	Status         bool            `json:"status"`
	StatusCode     int             `json:"statusCode"`
	Message        string          `json:"message"`
	Data           interface{}     `json:"data,omitempty"`
	PaginationData *PaginationData `json:"paginationData,omitempty"`
	Error          string          `json:"error,omitempty"`
}

// PaginationData represents pagination information
type PaginationData struct {
	Total           int64 `json:"total"`
	Page            int   `json:"page"`
	TotalPage       int   `json:"totalPage"`
	Limit           int   `json:"limit"`
	HasNextPage     bool  `json:"hasNextPage"`
	HasPreviousPage bool  `json:"hasPreviousPage"`
}

// ResponseBuilder helps build standardized responses
type ResponseBuilder struct {
	ctx *gin.Context
}

// NewResponseBuilder creates a new response builder
func NewResponseBuilder(c *gin.Context) *ResponseBuilder {
	return &ResponseBuilder{ctx: c}
}

// Success sends a successful response
func (r *ResponseBuilder) Success(message string, data interface{}) {
	response := StandardResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    message,
		Data:       data,
	}
	r.ctx.JSON(http.StatusOK, response)
}

// SuccessWithPagination sends a successful response with pagination data
func (r *ResponseBuilder) SuccessWithPagination(message string, data interface{}, total int64, page, limit int) {
	totalPages := int(math.Ceil(float64(total) / float64(limit)))
	
	paginationData := &PaginationData{
		Total:           total,
		Page:            page,
		TotalPage:       totalPages,
		Limit:           limit,
		HasNextPage:     page < totalPages,
		HasPreviousPage: page > 1,
	}

	response := StandardResponse{
		Status:         true,
		StatusCode:     http.StatusOK,
		Message:        message,
		Data:           data,
		PaginationData: paginationData,
	}
	r.ctx.JSON(http.StatusOK, response)
}

// Error sends an error response
func (r *ResponseBuilder) Error(statusCode int, message string, err ...string) {
	errorMsg := ""
	if len(err) > 0 {
		errorMsg = err[0]
	}

	response := StandardResponse{
		Status:     false,
		StatusCode: statusCode,
		Message:    message,
		Error:      errorMsg,
	}
	r.ctx.JSON(statusCode, response)
}

// BadRequest sends a 400 Bad Request response
func (r *ResponseBuilder) BadRequest(message string, err ...string) {
	r.Error(http.StatusBadRequest, message, err...)
}

// Unauthorized sends a 401 Unauthorized response
func (r *ResponseBuilder) Unauthorized(message string, err ...string) {
	r.Error(http.StatusUnauthorized, message, err...)
}

// NotFound sends a 404 Not Found response
func (r *ResponseBuilder) NotFound(message string, err ...string) {
	r.Error(http.StatusNotFound, message, err...)
}

// Conflict sends a 409 Conflict response
func (r *ResponseBuilder) Conflict(message string, err ...string) {
	r.Error(http.StatusConflict, message, err...)
}

// InternalServerError sends a 500 Internal Server Error response
func (r *ResponseBuilder) InternalServerError(message string, err ...string) {
	r.Error(http.StatusInternalServerError, message, err...)
}

// Created sends a 201 Created response
func (r *ResponseBuilder) Created(message string, data interface{}) {
	response := StandardResponse{
		Status:     true,
		StatusCode: http.StatusCreated,
		Message:    message,
		Data:       data,
	}
	r.ctx.JSON(http.StatusCreated, response)
}

// Updated sends a success response for updates
func (r *ResponseBuilder) Updated(message string, data ...interface{}) {
	var responseData interface{}
	if len(data) > 0 {
		responseData = data[0]
	}

	response := StandardResponse{
		Status:     true,
		StatusCode: http.StatusOK,
		Message:    message,
		Data:       responseData,
	}
	r.ctx.JSON(http.StatusOK, response)
}

// NoContent sends a 204 No Content response
func (r *ResponseBuilder) NoContent(message string) {
	response := StandardResponse{
		Status:     true,
		StatusCode: http.StatusNoContent,
		Message:    message,
	}
	r.ctx.JSON(http.StatusNoContent, response)
}