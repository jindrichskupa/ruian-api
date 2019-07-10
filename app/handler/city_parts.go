package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jindrichskupa/ruian-api/app/model"
	"github.com/jinzhu/gorm"
)

// GetAllCityParts from db
func GetAllCityParts(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	Cities := []model.CityPart{}
	db.Find(&Cities)
	respondJSON(w, http.StatusOK, Cities)
}

// GetCityPart by name
func GetCityPart(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	name = normalizeNameSearch(name)
	City := getCityPartOr404(db, name, w, r)
	if City == nil {
		return
	}
	respondJSON(w, http.StatusOK, City)
}

// getCityPartOr404 gets a City instance if exists, or respond the 404 error otherwise
func getCityPartOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *[]model.CityPart {
	CityParts := []model.CityPart{}
	err := db.Preload("City").Where("name_search LIKE ?", "%"+name+"%").Find(&CityParts).Error

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &CityParts
}
