package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zero-frost/xerospy-stats/app/model"
	"net/http"
)

type ActionTypeController struct {
	*Controller
}

func (atc *ActionTypeController) GetActionType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	actionType := model.ActionType{}
	atc.DB.Where("_id = ?", vars["action_type_id"]).First(&actionType)
	actionTypeData, err := json.Marshal(actionType)
	if err != nil {
		fmt.Fprint(w, "Error: No actionType found with id"+vars["action_type_id"])
		return
	}
	fmt.Fprint(w, string(actionTypeData))
}

func (atc *ActionTypeController) GetActionTypes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	actionTypes := make([]model.ActionType, 0)
	atc.DB.Where("_id = ?", vars["action_type_id"]).First(&actionTypes)
	actionTypeData, err := json.Marshal(actionTypes)
	if err != nil {
		fmt.Fprint(w, "Error: No actionType found with id"+vars["action_type_id"])
		return
	}
	fmt.Fprint(w, string(actionTypeData))
}

func (atc *ActionTypeController) UpdateActionType(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = atc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	vars := mux.Vars(r)

	var updateData model.ActionType

	if atc.DB.Where("_id = ?", vars["action_type_id"]).First(&updateData).Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	atc.DB.Save(updateData)
	fmt.Fprintf(w, "{\"status\": \"done\"}")
}

func (atc *ActionTypeController) AddActionType(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = atc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	var updateData model.ActionType

	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var oldData model.ActionType
	if atc.DB.Where("_id = ?", updateData.Id).First(&oldData).Error != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error: a actionType with ID `%s` already exists", oldData.Id)
	}

	atc.DB.Create(&updateData)
	fmt.Fprintf(w, "{\"status\": \"done\"}")

}

func (atc *ActionTypeController) DeleteActionType(w http.ResponseWriter, r *http.Request) {

	var err error
	var statusCode int
	if statusCode, _, err = atc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	vars := mux.Vars(r)

	var actionTypeData model.ActionType
	atc.DB.Where("_id = ?", vars["action_type_id"]).First(&actionTypeData)
	if err = atc.DB.Delete(actionTypeData).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
