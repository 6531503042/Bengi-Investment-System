package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidationError represents a single field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// ValidateStruct validates a struct and returns formatted errors
func ValidateStruct(s interface{}) *ValidationErrors {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var errors []ValidationError
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, ValidationError{
			Field:   toSnakeCase(err.Field()),
			Message: getErrorMessage(err),
		})
	}

	return &ValidationErrors{Errors: errors}
}

// ParseAndValidate parses request body and validates it
func ParseAndValidate(c *fiber.Ctx, out interface{}) *ValidationErrors {
	if err := c.BodyParser(out); err != nil {
		return &ValidationErrors{
			Errors: []ValidationError{
				{Field: "body", Message: "Invalid request body"},
			},
		}
	}
	return ValidateStruct(out)
}

// getErrorMessage returns a human-readable error message
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Must be at least " + err.Param() + " characters"
	case "max":
		return "Must be at most " + err.Param() + " characters"
	case "len":
		return "Must be exactly " + err.Param() + " characters"
	case "eqfield":
		return "Must match " + toSnakeCase(err.Param())
	case "gt":
		return "Must be greater than " + err.Param()
	case "gte":
		return "Must be greater than or equal to " + err.Param()
	case "lt":
		return "Must be less than " + err.Param()
	case "lte":
		return "Must be less than or equal to " + err.Param()
	case "oneof":
		return "Must be one of: " + err.Param()
	case "alphanum":
		return "Must contain only alphanumeric characters"
	default:
		return "Invalid value"
	}
}

// toSnakeCase converts PascalCase/camelCase to snake_case
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
