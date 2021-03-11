package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jindrichskupa/ruian-api/app/model"
	"github.com/jinzhu/gorm"
)

// GetAllPlaces from db
func GetAllPlaces(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	Places := []model.Place{}
	db.Find(&Places)
	respondJSON(w, http.StatusOK, Places)
}

// GetPlace by id
func GetPlace(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
	}

	Place := getPlaceOr404(db, id, w, r)
	if Place == nil {
		return
	}
	respondJSON(w, http.StatusOK, Place)
}

// SearchPlace by filter
func SearchPlace(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	Place := searchPlaceOr404(db, w, r)
	if Place == nil {
		return
	}
	respondJSON(w, http.StatusOK, Place)
}

// getCityOr404 gets a City instance if exists, or respond the 404 error otherwise
func getPlaceOr404(db *gorm.DB, id int, w http.ResponseWriter, r *http.Request) *[]model.Place {

	Places := []model.Place{}
	err := db.Preload("Street").Preload("City").Preload("CityPart").Find(&Places, id).Error

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &Places
}

// getCityOr404 gets a City instance if exists, or respond the 404 error otherwise
func searchPlaceOr404(db *gorm.DB, w http.ResponseWriter, r *http.Request) *[]model.Place {
	queryValues := r.URL.Query()

	//log.Println(normalizeNameSearch("Školní 105, Zruč - Senec, 33008"))

	log.Println(queryValues)
	filters := []string{"street", "city_part", "city"}
	numbers := []string{"p", "e"}
	likes := []string{"zip", "o"}

	tables := map[string]string{
		"place":     db.NewScope(&model.Place{}).TableName(),
		"street":    db.NewScope(&model.Street{}).TableName(),
		"city_part": db.NewScope(&model.CityPart{}).TableName(),
		"city":      db.NewScope(&model.City{}).TableName(),
	}

	tx := db.Where("1 = 1")

	for _, filter := range filters {
		if queryValues[filter] != nil {
			value := normalizeNameSearch(queryValues[filter][0])
			value = "%" + value + "%"
			log.Println("Key [", filter, "]: ", value)
			join := fmt.Sprintf("INNER JOIN %s ON %s.%s_id = %s.id AND %s.name_search LIKE ?",
				tables[filter],
				tables["place"],
				filter,
				tables[filter],
				tables[filter],
			)
			log.Println("Key [", filter, "]: ", join, value)
			tx = tx.Joins(join, value)
		}
	}

	if queryValues["lat"] != nil && queryValues["lng"] != nil {
		gpsRange := 100
		if queryValues["range"] != nil {
			gpsRange, err := strconv.Atoi(queryValues["range"][0])
			if err != nil || gpsRange > 1000 {
				gpsRange = 100
			}
		}
		// geoFilter := fmt.Sprintf("earth_box(ll_to_earth(latitude,longitude), %s) @> ll_to_earth(%s,%s)",
		// 	queryValues["range"][0],
		// 	queryValues["latitude"][0],
		// 	queryValues["longitude"][0],
		// )
		geoFilter := fmt.Sprintf("point(latitude,longitude) <@> point(%s,%s) < %d::float8/1609::float8",
			queryValues["lat"][0],
			queryValues["lng"][0],
			gpsRange,
		)

		tx = tx.Where(geoFilter)
	}

	for _, number := range numbers {
		if queryValues[number] != nil {
			value := queryValues[number][0]
			log.Println("Key [", number, "]: ", value)
			tx = tx.Where(number+" = ?", value)
		}
	}

	for _, like := range likes {
		if queryValues[like] != nil {
			value := queryValues[like][0]
			log.Println("Key [", like, "]: ", value)
			tx = tx.Where(like+" LIKE ?", "%"+value+"%")
		}
	}

	if queryValues["limit"] != nil {
		limit, err := strconv.Atoi(queryValues["limit"][0])
		if err != nil {
			limit = 200
		}
		tx = tx.Limit(limit)
	} else {
		tx = tx.Limit(200)
	}

	Places := []model.Place{}
	err := tx.Preload("Street").
		Preload("City").
		Preload("CityPart").
		Select("DISTINCT \"view_address_places\".*").
		Find(&Places).Error

	if err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &Places
}
