package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jindrichskupa/ruian-api/app/model"
	"github.com/jinzhu/gorm"
)

// GetAllCities from db
func GetAllCities(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	Cities := []model.City{}
	db.Find(&Cities)
	respondJSON(w, http.StatusOK, Cities)
}

// GetCity by name
func GetCity(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	name = normalizeNameSearch(name)
	Cities := getCityOr404(db, name, w, r)
	if Cities == nil {
		return
	}
	respondJSON(w, http.StatusOK, Cities)
}

// getCityOr404 gets a City instance if exists, or respond the 404 error otherwise
func getCityOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *[]model.City {
	Cities := []model.City{}
	err := db.Where("name_search LIKE ?", "%"+name+"%").Find(&Cities).Error

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &Cities
}
