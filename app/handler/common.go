package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// respondJSON makes the response with payload as json format
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// respondError makes the error response with payload as json format
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

func standardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func normalizeNameSearch(name string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)
	ascii, _, _ := transform.String(t, name)
	ascii = strings.ToLower(ascii)

	reg, err := regexp.Compile("[^a-z0-9 ]+")
	if err != nil {
		log.Println(err)
	}
	ascii = reg.ReplaceAllString(ascii, " ")

	return standardizeSpaces(ascii)
}
