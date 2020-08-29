package lib

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const NoRowFound = "record not found"

// NewInternalServiceErr returns internal err
func NewInternalServiceErr(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// NewBadRequestErr returns bad request err
func NewBadRequestErr(field string) gin.H {
	return gin.H{"error": fmt.Sprintf("%s field is invalid", field)}
}

// NewNotFoundErr returns not found err
func NewNotFoundErr(resource string, params interface{}) gin.H {
	return gin.H{"error": fmt.Sprintf("%s(#%v) is not found", resource, params)}
}

// NewConflict returns conflict err
func NewConflict() gin.H {
	return gin.H{"error": "resource is already existed"}
}

// NewUnauthorized return unauthorized err
func NewUnauthorized() gin.H {
	return gin.H{"error": "user is not authorized"}
}
