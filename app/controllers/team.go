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
	fmt.Printf("%+v\n", team)
	team_data, err := json.Marshal(team)
	if err != nil {
		fmt.Fprint(w, "Error: No team found with id"+vars["team_id"])
		return
	}
	fmt.Fprint(w, string(team_data))
}

func (tc *TeamController) GetTeams(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teams := make([]model.Team, 0)
	tc.DB.Where("_id = ?", vars["team_id"]).First(&teams)
	team_data, err := json.Marshal(teams)
	if err != nil {
		fmt.Fprint(w, "Error: No team found with id"+vars["team_id"])
		return
	}
	fmt.Fprint(w, string(team_data))
}
