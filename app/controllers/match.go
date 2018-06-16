package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zero-frost/xerospy-stats/app/model"
	"net/http"
)

type MatchController struct {
	*Controller
}

func (mc *MatchController) GetMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	match := model.Match{}
	mc.DB.Where("_id = ?", vars["match_id"]).First(&match)
	matchData, err := json.Marshal(match)
	if err != nil {
		fmt.Fprint(w, "Error: No match found with id"+vars["match_id"])
		return
	}
	fmt.Fprint(w, string(matchData))
}

func (mc *MatchController) GetMatches(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	matches := make([]model.Match, 0)
	mc.DB.Where("_id = ?", vars["match_id"]).First(&matches)
	matchData, err := json.Marshal(matches)
	if err != nil {
		fmt.Fprint(w, "Error: No match found with id"+vars["match_id"])
		return
	}
	fmt.Fprint(w, string(matchData))
}

func (mc *MatchController) UpdateMatch(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = mc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	vars := mux.Vars(r)

	var updateData model.Match

	if mc.DB.Where("_id = ?", vars["match_id"]).First(&updateData).Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	mc.DB.Save(updateData)
	fmt.Fprintf(w, "{\"status\": \"done\"}")
}

func (mc *MatchController) AddMatch(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = mc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	var updateData model.Match

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var oldData model.Match
	if mc.DB.Where("_id = ?", updateData.Id).First(&oldData).Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: a match with ID `%s` already exists", oldData.Id)
	}

	mc.DB.Create(&updateData)
	fmt.Fprintf(w, "{\"status\": \"done\"}")

}

func (mc *MatchController) DeleteMatch(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = mc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	vars := mux.Vars(r)

	var matchData model.Match
	mc.DB.Where("_id = ?", vars["match_id"]).First(&matchData)
	if err = mc.DB.Delete(matchData).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
