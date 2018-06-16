package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zero-frost/xerospy-stats/app/model"
	"net/http"
)

type TeamController struct {
	*Controller
}

func (tc *TeamController) GetTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	team := model.Team{}
	tc.DB.Where("_id = ?", vars["team_id"]).First(&team)
	teamData, err := json.Marshal(team)
	if err != nil {
		fmt.Fprint(w, "Error: No team found with id"+vars["team_id"])
		return
	}
	fmt.Fprint(w, string(teamData))
}

func (tc *TeamController) GetTeams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teams := make([]model.Team, 0)
	tc.DB.Where("_id = ?", vars["team_id"]).First(&teams)
	teamData, err := json.Marshal(teams)
	if err != nil {
		fmt.Fprint(w, "Error: No team found with id"+vars["team_id"])
		return
	}
	fmt.Fprint(w, string(teamData))
}

func (tc *TeamController) UpdateTeam(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = tc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	vars := mux.Vars(r)

	var updateData model.Team

	tc.DB.Where("_id = ?", vars["team_id"]).First(&updateData)

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tc.DB.Save(updateData)
	fmt.Fprintf(w, "{\"status\": \"done\"}")
}

func (tc *TeamController) AddTeam(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = tc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	var updateData model.Team

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var oldData model.Team
	tc.DB.Where("_id = ?", updateData.Id).First(&oldData)

	if oldData.Id == 0 && oldData.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: a team with ID `%s` already exists", oldData.Id)
	}

	tc.DB.Create(&updateData)
	fmt.Fprintf(w, "{\"status\": \"done\"}")

}

func (tc *TeamController) DeleteTeam(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = tc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	vars := mux.Vars(r)

	var teamData model.Team
	tc.DB.Where("_id = ?", vars["team_id"]).First(&teamData)
	if err = tc.DB.Delete(teamData).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
