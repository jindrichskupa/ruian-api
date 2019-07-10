package handler

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jindrichskupa/ruian-api/app/model"
	"github.com/jinzhu/gorm"
)

// GetAllCadastralTerritories from db
func GetAllCadastralTerritories(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	CadastralTerritories := []model.CadastralTerritory{}
	db.Preload("City").Find(&CadastralTerritories)
	respondJSON(w, http.StatusOK, CadastralTerritories)
}

// GetCadastralTerritory by name
func GetCadastralTerritory(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	name := vars["name"]
	name = normalizeNameSearch(name)
	CadastralTerritory := getCadastralTerritoryOr404(db, name, w, r)
	if CadastralTerritory == nil {
		return
	}
	respondJSON(w, http.StatusOK, CadastralTerritory)
}

// getCadastralTerritoryOr404 gets a CadastralTerritory instance if exists, or respond the 404 error otherwise
func getCadastralTerritoryOr404(db *gorm.DB, name string, w http.ResponseWriter, r *http.Request) *[]model.CadastralTerritory {
	CadastralTerritories := []model.CadastralTerritory{}
	err := db.Preload("City").Where("name_search LIKE ?", "%"+name+"%").Find(&CadastralTerritories).Error

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &CadastralTerritories
}
