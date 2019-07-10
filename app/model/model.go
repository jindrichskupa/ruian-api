package model

import (
	"encoding/json"
	"fmt"

	"github.com/jinzhu/gorm"
	// import PostgreSQL dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// City structure
type City struct {
	ID         uint   `gorm:"primary_key" json:"id"`
	Name       string `json:"name"`
	NameSearch string `json:"name_search"`
	CityParts  []CityPart
	Streets    []Street
	Places     []Place
}

// ToString convert city struct to string
func (c City) ToString() string {
	return c.Name
}

// CityJSON structure
type CityJSON struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// MarshalJSON converts City to CityJSON
func (c City) MarshalJSON() ([]byte, error) {
	m := CityJSON{
		ID:   c.ID,
		Name: c.Name,
	}
	return json.Marshal(m)
}

// CityPart structure
type CityPart struct {
	ID         uint   `gorm:"primary_key" json:"id"`
	Name       string `json:"name"`
	NameSearch string `json:"name_search"`
	CityID     uint   `json:"city_id"`
	City       City   `gorm:"foreignkey:CityID;association_foreignkey:ID" json:"city"`
	Streets    []Street
	Places     []Place
}

// ToString convert city part struct to string
func (cp CityPart) ToString() string {
	return cp.Name
}

// CityPartJSON structure
type CityPartJSON struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

// MarshalJSON converts CityPart to CityPartJSON
func (cp CityPart) MarshalJSON() ([]byte, error) {
	m := CityPartJSON{
		ID:   cp.ID,
		Name: cp.Name,
		City: cp.City.ToString(),
	}
	return json.Marshal(m)
}

// Street structure
type Street struct {
	ID         uint     `gorm:"primary_key" json:"id"`
	Name       string   `json:"name"`
	NameSearch string   `json:"name_search"`
	CityID     uint     `json:"city_id"`
	CityPartID uint     `json:"city_part_id"`
	City       City     `gorm:"foreignkey:CityID;association_foreignkey:ID" json:"city"`
	CityPart   CityPart `gorm:"foreignkey:CityPartID;association_foreignkey:ID" json:"city_part"`
	Places     []Place
}

// ToString convert street struct to string
func (s Street) ToString() string {
	return s.Name
}

// StreetJSON structure
type StreetJSON struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	City     string `json:"city"`
	CityPart string `json:"city_part"`
}

// MarshalJSON converts Street to StreetJSON
func (s Street) MarshalJSON() ([]byte, error) {
	m := StreetJSON{
		ID:       s.ID,
		Name:     s.Name,
		City:     s.City.ToString(),
		CityPart: s.CityPart.ToString(),
	}
	return json.Marshal(m)
}

// Place structure
type Place struct {
	ID         uint     `gorm:"primary_key" json:"id"`
	E          uint32   `json:"e,omitempty"`
	P          uint32   `json:"p,omitempty"`
	O          string   `gorm:"size:10" json:"o,omitempty"`
	Zip        string   `gorm:"size:5" json:"zip"`
	X          float64  `json:"x,omitempty"`
	Y          float64  `json:"y,omitempty"`
	Longitude  float64  `json:"lng,omitempty"`
	Latitude   float64  `json:"lat,omitempty"`
	StreetID   uint     `json:"-"`
	CityID     uint     `json:"-"`
	CityPartID uint     `json:"-"`
	Street     Street   `gorm:"foreignkey:StreetID;association_foreignkey:ID" json:"street"`
	City       City     `gorm:"foreignkey:CityID;association_foreignkey:ID" json:"city"`
	CityPart   CityPart `gorm:"foreignkey:CityPartID;association_foreignkey:ID" json:"city_part"`
}

// ToString convert city part struct to string
func (p Place) ToString() string {
	return fmt.Sprintf("%s %s, %s, %s, %s",
		p.Street.ToString(),
		p.PlaceNumber(),
		p.CityPart.ToString(),
		p.City.ToString(),
		p.Zip,
	)
}

// PlaceNumber returns correct address place number
func (p Place) PlaceNumber() string {
	number := ""
	if p.P == 0 {
		number = fmt.Sprintf("E%d", p.E)
	} else {
		if p.O == "" {
			number = fmt.Sprintf("%d", p.P)
		} else {
			number = fmt.Sprintf("%d/%s", p.P, p.O)
		}
	}
	return number
}

// PlaceJSON Place representation for JSON
type PlaceJSON struct {
	ID            uint    `json:"id"`
	E             uint32  `json:"e,omitempty"`
	P             uint32  `json:"p,omitempty"`
	O             string  `json:"o,omitempty"`
	Zip           string  `json:"zip"`
	Longitude     float64 `json:"lng"`
	Latitude      float64 `json:"lat"`
	Street        string  `json:"street"`
	City          string  `json:"city"`
	CityPart      string  `json:"city_part"`
	AddressString string  `json:"address_string"`
}

// MarshalJSON converts Place to PlaceJSON
func (p Place) MarshalJSON() ([]byte, error) {
	m := PlaceJSON{
		ID:            p.ID,
		E:             p.E,
		P:             p.P,
		O:             p.O,
		Zip:           p.Zip,
		Longitude:     p.Longitude,
		Latitude:      p.Latitude,
		Street:        p.Street.ToString(),
		City:          p.City.ToString(),
		CityPart:      p.CityPart.ToString(),
		AddressString: p.ToString(),
	}
	return json.Marshal(m)
}

// CadastralTerritory structure
type CadastralTerritory struct {
	ID         uint   `gorm:"primary_key" json:"id"`
	Name       string `json:"name"`
	NameSearch string `json:"name_search"`
	CityID     uint   `json:"city_id"`
	City       City   `gorm:"foreignkey:CityID;association_foreignkey:ID" json:"city"`
}

// CadastralTerritoryJSON CadastralTerritory representation for JSON
type CadastralTerritoryJSON struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

// MarshalJSON converts CadastralTerritory to CadastralTerritoryJSON
func (ct CadastralTerritory) MarshalJSON() ([]byte, error) {
	m := CadastralTerritoryJSON{
		ID:   ct.ID,
		Name: ct.Name,
		City: ct.City.ToString(),
	}
	return json.Marshal(m)
}

// DBMigrate will create and migrate the tables, and then make the some relationships if necessary
func DBMigrate(db *gorm.DB) *gorm.DB {
	return db
}
