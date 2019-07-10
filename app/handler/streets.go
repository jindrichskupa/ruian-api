package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jindrichskupa/ruian-api/app/model"
	"github.com/jinzhu/gorm"
)

// GetAllStreets from db
func GetAllStreets(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	Cities := []model.Street{}
	db.Find(&Cities)
	respondJSON(w, http.StatusOK, Cities)
}

// GetStreet by name
func GetStreet(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	name = normalizeNameSearch(name)
	Street := getStreetOr404(db, name, w, r)
	if Street == nil {
		return
	}
	respondJSON(w, http.StatusOK, Street)
}

// getStreetOr404 gets a Streets instance if exists, or respond the 404 error otherwise
func getStreetOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *[]model.Street {
	Streets := []model.Street{}
	err := db.Preload("City").Preload("CityPart").Where("name_search LIKE ?", "%"+name+"%").Find(&Streets).Error

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &Streets
}
