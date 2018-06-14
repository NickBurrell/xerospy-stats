package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zero-frost/xerospy-stats/app/model"
	"net/http"
	"reflect"
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

func (tc *TeamController) UpdateTeam(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	c, err := r.Cookie("session_token")

	if err != nil {
		if err == http.ErrNoCookie {
			return
		}
		return
	}

	sessionToken := c.Value

	var updateData model.Team

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var res string
	if res, err = tc.Cache.Get(sessionToken).Result(); res == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var oldData model.Team
	tc.DB.Where("_id = ?", vars["team_id"]).First(&oldData)

	current := reflect.ValueOf(oldData).Elem()
	update := reflect.ValueOf(updateData).Elem()
	for i := 0; i < update.NumField(); i++ {
		currentField := current.Field(i)
		updateFieldValue := reflect.Value(update.Field(i))
		if updateFieldValue.String() != "" {
			currentField.Set(updateFieldValue)
		}
	}

	tc.DB.Save(updateData)
	fmt.Fprintf(w, "{\"status\": \"done\"}")
}

func (tc *TeamController) AddTeam(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session_token")

	if err != nil {
		if err == http.ErrNoCookie {
			return
		}
		return
	}

	sessionToken := c.Value

	var updateData model.Team

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var res string
	if res, err = tc.Cache.Get(sessionToken).Result(); res == "" {
		w.WriteHeader(http.StatusUnauthorized)
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

	c, err := r.Cookie("session_token")

	if err != nil {
		if err == http.ErrNoCookie {
			return
		}
		return
	}

	sessionToken := c.Value

	var updateData model.Team

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var res string
	if res, err = tc.Cache.Get(sessionToken).Result(); res == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

}
