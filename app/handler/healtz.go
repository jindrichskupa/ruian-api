package handler

import (
	"net/http"

	"github.com/jinzhu/gorm"
)

// HealthCheck struct
type HealthCheck struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// GetHealtStatus retuns application status info
func GetHealtStatus(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	health := HealthCheck{Status: "OK", Message: "I'm alive"}
	respondJSON(w, http.StatusOK, health)
}
