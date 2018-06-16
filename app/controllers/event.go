package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zero-frost/xerospy-stats/app/model"
	"net/http"
)

type EventController struct {
	*Controller
}

func (ec *EventController) GetEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	event := model.Event{}
	ec.DB.Where("_id = ?", vars["event_id"]).First(&event)
	fmt.Printf("%+v\n", event)
	eventData, err := json.Marshal(event)
	if err != nil {
		fmt.Fprint(w, "Error: No team found with id"+vars["event_id"])
		return
	}
	fmt.Fprint(w, string(eventData))
}

func (ec *EventController) GetEvents(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	teams := make([]model.Event, 0)
	ec.DB.Where("_id = ?", vars["event_id"]).First(&teams)
	eventData, err := json.Marshal(teams)
	if err != nil {
		fmt.Fprint(w, "Error: No team found with id"+vars["event_id"])
		return
	}
	fmt.Fprint(w, string(eventData))
}

func (ec *EventController) UpdateEvent(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = ec.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	var updateData model.Event

	vars := mux.Vars(r)

	ec.DB.Where("_id = ?", vars["event_id"]).First(&updateData)

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ec.DB.Save(updateData)
	fmt.Fprintf(w, "{\"status\": \"done\"}")
}

func (ec *EventController) AddEvent(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = ec.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	var updateData model.Event

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var oldData model.Team
	if err := ec.DB.Where("_id = ?", updateData.Id).First(&oldData).Error; err == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: a team with ID `%s` already exists", oldData.Id)
	}

	ec.DB.Create(&updateData)
	fmt.Fprintf(w, "{\"status\": \"done\"}")

}

func (ec *EventController) DeleteEvent(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = ec.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	vars := mux.Vars(r)

	var eventData model.Event
	ec.DB.Where("_id = ?", vars["event_id"]).First(&eventData)
	if err = ec.DB.Delete(eventData).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
