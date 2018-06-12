package controllers

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"github.com/zero-frost/xerospy-stats/app"
	"github.com/zero-frost/xerospy-stats/app/model"
	"net/http"
)

type LoginController struct {
	Db *gorm.Db
}

type loginData struct {
	Username string
	Password string
	Email    string
}

func (lc LoginController) Login(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	var credentials loginData

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !model.ValidateEmail(credentials.Email) {
		fmt.Fprint("Invalid Email %s\n", credentials.Email)
		return
	}

	var userData model.User
	lc.Db.Where("Username = ?", credentials.Username).First(&userData)

	if userData.Username == "" {
		fmt.Fprint("No user found by the username %s\n", credentials.Username)
		return
	}
	if []byte(userData.Password) != Sum256([]byte(credentials.Password+Salt)) {
		fmt.Fprint("Incorrect Password\n")
		return
	}

	sessionToken := uuid.NewV4().String()
}
