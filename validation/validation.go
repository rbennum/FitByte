package validation

import (
	"github.com/go-playground/validator/v10"
	"net/url"
)

// IsValidURI is a custom validation function
func IsValidURI(fl validator.FieldLevel) bool {
	uriStr, ok := fl.Field().Interface().(string)
	if !ok {
		// If the field is not a string (e.g., nil pointer), skip validation.
		return true
	}

	if uriStr == "" {
		// Allow empty URIs with `omitempty`.
		return true
	}

	// Parse the URI using net/url
	parsedURI, err := url.Parse(uriStr)
	if err != nil || parsedURI.Scheme == "" || parsedURI.Host == "" {
		// Invalid URI, must have scheme (e.g., "http") and host (e.g., "example.com")
		return false
	}

	// Ensure there's a path (e.g., "/image.jpg")
	return parsedURI.Path != "" && parsedURI.Path != "/"
}
