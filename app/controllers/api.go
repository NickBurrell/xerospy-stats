package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/zero-frost/xerospy-stats/app"
	"github.com/zero-frost/xerospy-stats/app/model"
	"net/http"
)

type TeamController struct {
	db *gorm.DB
}

func GetTeam(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	team := model.Team{}
	app.GetDatabase().Where("_id = ?", vars["team_id"]).First(&team)
	team_data, err := json.Marshal(team)
	if err != nil {
		fmt.Fprint(w, "Error: No team found with id"+vars["team_id"])
		return
	}
	fmt.Fprint(w, string(team_data))
}
