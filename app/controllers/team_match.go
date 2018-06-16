package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zero-frost/xerospy-stats/app/model"
	"net/http"
)

type TeamMatchController struct {
	*Controller
}

func (tmc *TeamMatchController) GetTeamMatch(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamMatch := model.TeamMatch{}
	tmc.DB.Where("_id = ?", vars["teamMatch_id"]).First(&teamMatch)
	teamMatchData, err := json.Marshal(teamMatch)
	if err != nil {
		fmt.Fprint(w, "Error: No teamMatch found with id"+vars["team_match_id"])
		return
	}
	fmt.Fprint(w, string(teamMatchData))
}

func (tmc *TeamMatchController) GetTeamMatches(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamMatches := make([]model.TeamMatch, 0)
	tmc.DB.Where("_id = ?", vars["teamMatch_id"]).First(&teamMatches)
	teamMatchData, err := json.Marshal(teamMatches)
	if err != nil {
		fmt.Fprint(w, "Error: No teamMatch found with id"+vars["team_match_id"])
		return
	}
	fmt.Fprint(w, string(teamMatchData))
}

func (tmc *TeamMatchController) UpdateTeamMatch(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = tmc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	vars := mux.Vars(r)

	var updateData model.TeamMatch

	tmc.DB.Where("_id = ?", vars["teamMatch_id"]).First(&updateData)

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tmc.DB.Save(updateData)
	fmt.Fprintf(w, "{\"status\": \"done\"}")
}

func (tmc *TeamMatchController) AddTeamMatch(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = tmc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	var updateData model.TeamMatch

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var oldData model.TeamMatch
	if tmc.DB.Where("_id = ?", updateData.Id).First(&oldData).Error == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: a teamMatch with ID `%s` already exists", oldData.Id)
	}

	tmc.DB.Create(&updateData)
	fmt.Fprintf(w, "{\"status\": \"done\"}")

}

func (tmc *TeamMatchController) DeleteTeamMatch(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = tmc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	vars := mux.Vars(r)

	var teamMatchData model.TeamMatch
	tmc.DB.Where("_id = ?", vars["teamMatch_id"]).First(&teamMatchData)
	if err = tmc.DB.Delete(teamMatchData).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
