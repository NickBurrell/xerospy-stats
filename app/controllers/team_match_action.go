package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zero-frost/xerospy-stats/app/model"
	"net/http"
)

type TeamMatchActionController struct {
	*Controller
}

func (tmac *TeamMatchActionController) GetTeamMatchAction(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamMatchAction := model.TeamMatchAction{}
	tmac.DB.Where("_id = ?", vars["team_match_action_id"]).First(&teamMatchAction)
	teamMatchActionData, err := json.Marshal(teamMatchAction)
	if err != nil {
		fmt.Fprint(w, "Error: No teamMatchAction found with id"+vars["team_match_action_id"])
		return
	}
	fmt.Fprint(w, string(teamMatchActionData))
}

func (tmac *TeamMatchActionController) GetTeamMatchActions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teamMatchActions := make([]model.TeamMatchAction, 0)
	tmac.DB.Where("_id = ?", vars["team_match_action_id"]).First(&teamMatchActions)
	teamMatchActionData, err := json.Marshal(teamMatchActions)
	if err != nil {
		fmt.Fprint(w, "Error: No teamMatchAction found with id"+vars["team_match_action_id"])
		return
	}
	fmt.Fprint(w, string(teamMatchActionData))
}

func (tmac *TeamMatchActionController) UpdateTeamMatchAction(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = tmac.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	vars := mux.Vars(r)

	var updateData model.TeamMatchAction

	if tmac.DB.Where("_id = ?", vars["team_match_action_id"]).First(&updateData).Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tmac.DB.Save(updateData)
	fmt.Fprintf(w, "{\"status\": \"done\"}")
}

func (tmac *TeamMatchActionController) AddTeamMatchAction(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = tmac.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	var updateData model.TeamMatchAction

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var oldData model.TeamMatchAction
	if tmac.DB.Where("_id = ?", updateData.Id).First(&oldData).Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: a teamMatchAction with ID `%s` already exists", oldData.Id)
	}

	tmac.DB.Create(&updateData)
	fmt.Fprintf(w, "{\"status\": \"done\"}")

}

func (tmac *TeamMatchActionController) DeleteTeamMatchAction(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = tmac.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	vars := mux.Vars(r)

	var teamMatchActionData model.TeamMatchAction

	tmac.DB.Where("_id = ?", vars["team_match_action_id"]).First(&teamMatchActionData)
	if err = tmac.DB.Delete(teamMatchActionData).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
