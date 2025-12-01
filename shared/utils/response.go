package response

import "github.com/gin-gonic/gin"

// Meta response metadata according to API contract
type Meta struct {
	Code    int    `json:"code"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// APIResponse standard API response structure with meta wrapper
type APIResponse struct {
	Meta   Meta                `json:"meta"`
	Data   interface{}         `json:"data"`
	Errors map[string][]string `json:"errors,omitempty"`
}

// SuccessResponse sends success response
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Meta: Meta{
			Code:    statusCode,
			Status:  "success",
			Message: message,
		},
		Data: data,
	})
}

// ErrorResponse sends error response
func ErrorResponse(c *gin.Context, statusCode int, message string, err error) {
	c.JSON(statusCode, APIResponse{
		Meta: Meta{
			Code:    statusCode,
			Status:  "error",
			Message: message,
		},
		Data: "",
	})
}

// ValidationErrorResponse sends validation error response
func ValidationErrorResponse(c *gin.Context, errors map[string][]string) {
	c.JSON(400, APIResponse{
		Meta: Meta{
			Code:    400,
			Status:  "error",
			Message: "Validation error",
		},
		Data:   "",
		Errors: errors,
	})
}
