package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jindrichskupa/ruian-api/app/handler"
	"github.com/jindrichskupa/ruian-api/app/model"
	"github.com/jindrichskupa/ruian-api/config"
	"github.com/jinzhu/gorm"
)

// App has router and db instances
type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

// Initialize application with predefined configuration
func (a *App) Initialize(config *config.Config) {
	dbURI := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable",
		config.DB.Dialect,
		config.DB.Username,
		config.DB.Password,
		config.DB.Hostname,
		config.DB.Port,
		config.DB.Name)

	db, err := gorm.Open(config.DB.Dialect, dbURI)
	for tries := 0; tries < 10 && err != nil; tries++ {
		time.Sleep(5 * time.Second)
		log.Println("Trying to reconnect (", tries, ") db: ", err)
		db, err = gorm.Open(config.DB.Dialect, dbURI)
	}
	if err != nil {
		log.Fatal("Could not connect database: ", err)
	}
	db.LogMode(true)

	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return config.DB.Prefix + defaultTableName
	}

	a.DB = model.DBMigrate(db)
	a.Router = mux.NewRouter()
	a.setRouters()
	log.Println("RUIAN API application started")
}

// Set all required routers
func (a *App) setRouters() {
	// Routing for handling the projects
	a.Get("/healtz", a.GetHealtStatus)
	a.Get("/cities", a.GetAllCities)
	a.Get("/cities/{name}", a.GetCity)
	a.Get("/city_parts", a.GetAllCityParts)
	a.Get("/city_parts/{name}", a.GetCityPart)
	a.Get("/streets", a.GetAllStreets)
	a.Get("/streets/{name}", a.GetStreet)
	a.Get("/places", a.GetAllPlaces)
	a.Get("/places/search", a.SearchPlace)
	a.Get("/places/{id}", a.GetPlace)
	a.Get("/cadastral_territories", a.GetAllCadastralTerritories)
	a.Get("/cadastral_territories/{name}", a.GetCadastralTerritory)
}

// Get wraps the router for GET method
func (a *App) Get(path string, f func(w http.ResponseWriter, r *http.Request)) {
	a.Router.HandleFunc(path, f).Methods("GET")
}

// GetHealtStatus retuns application status info
func (a *App) GetHealtStatus(w http.ResponseWriter, r *http.Request) {
	handler.GetHealtStatus(a.DB, w, r)
}

// GetAllCities handlers to manage City Data
func (a *App) GetAllCities(w http.ResponseWriter, r *http.Request) {
	handler.GetAllCities(a.DB, w, r)
}

// GetCity handlers to manage City Data
func (a *App) GetCity(w http.ResponseWriter, r *http.Request) {
	handler.GetCity(a.DB, w, r)
}

// GetAllCityParts handlers to manage City Part Data
func (a *App) GetAllCityParts(w http.ResponseWriter, r *http.Request) {
	handler.GetAllCityParts(a.DB, w, r)
}

// GetCityPart handlers to manage City Part Data
func (a *App) GetCityPart(w http.ResponseWriter, r *http.Request) {
	handler.GetCityPart(a.DB, w, r)
}

// GetAllStreets handlers to manage Street Data
func (a *App) GetAllStreets(w http.ResponseWriter, r *http.Request) {
	handler.GetAllStreets(a.DB, w, r)
}

// GetStreet handlers to manage Street Data
func (a *App) GetStreet(w http.ResponseWriter, r *http.Request) {
	handler.GetStreet(a.DB, w, r)
}

// GetAllPlaces handlers to manage Place Data
func (a *App) GetAllPlaces(w http.ResponseWriter, r *http.Request) {
	handler.GetAllPlaces(a.DB, w, r)
}

// GetPlace handlers to manage Place Data
func (a *App) GetPlace(w http.ResponseWriter, r *http.Request) {
	handler.GetPlace(a.DB, w, r)
}

// SearchPlace handler to get places by filter
func (a *App) SearchPlace(w http.ResponseWriter, r *http.Request) {
	handler.SearchPlace(a.DB, w, r)
}

// GetAllCadastralTerritories handlers to manage Cadastral Territory Data
func (a *App) GetAllCadastralTerritories(w http.ResponseWriter, r *http.Request) {
	handler.GetAllCadastralTerritories(a.DB, w, r)
}

// GetCadastralTerritory handlers to manage Cadastral Territory Data
func (a *App) GetCadastralTerritory(w http.ResponseWriter, r *http.Request) {
	handler.GetCadastralTerritory(a.DB, w, r)
}

// Run the app on it's router
func (a *App) Run(host string) {
	log.Fatal(http.ListenAndServe(host, a.Router))
}
