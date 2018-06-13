package controllers

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"github.com/zero-frost/xerospy-stats/app"
	"github.com/zero-frost/xerospy-stats/app/model"
	"log"
	"net/http"
	"regexp"
	"time"
)

type LoginController struct {
	DB *gorm.DB
}

type loginData struct {
	Username string
	Password string
	Email    string
}

func (lc LoginController) Login(w http.ResponseWriter, r *http.Request) {
	var credentials loginData

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !model.ValidateEmail(credentials.Email) {
		fmt.Fprintf(w, "Invalid Email %s\n", credentials.Email)
		return
	}

	var userData model.User
	lc.DB.Where("username = ?", credentials.Username).First(&userData)

	if userData.Username == "" {
		fmt.Fprintf(w, "No user found by the username %s\n", credentials.Username)
		return
	}

	hashedPass := sha256.Sum256([]byte(credentials.Password + app.GetServerConfig().ServerSettings.Salt))

	if !bytes.Equal([]byte(userData.Password), hashedPass[:]) {
		fmt.Fprint(w, "Incorrect Password\n")
		return
	}

	sessionToken, _ := uuid.NewV4()

	err = app.GetCache().Set(sessionToken.String(), credentials.Username, time.Duration(300*time.Second)).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Path:    "/",
		Expires: time.Now().Add(300 * time.Second),
	})

	w.WriteHeader(http.StatusOK)

}

func (lc LoginController) Refresh(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := c.Value

	response, err := app.GetCache().Get(sessionToken).Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if response == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// At this point we know the user is valid

	newSessionToken, _ := uuid.NewV4()
	err = app.GetCache().Set(newSessionToken.String(), fmt.Sprint(response), time.Duration(300*time.Second)).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = app.GetCache().Del(sessionToken).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken.String(),
		Path:    "/",
		Expires: time.Now().Add(300 * time.Second),
	})
}

func (lc LoginController) Logout(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session_token")

	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := c.Value

	err = app.GetCache().Del(sessionToken).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.MaxAge = -1
	c.Value = ""
	c.Path = "/"
	http.SetCookie(w, c)
}

func (lc LoginController) CreateUser(w http.ResponseWriter, r *http.Request) {

	var credentials loginData

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if credentials.Email == "" || credentials.Password == "" || credentials.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var placeholderUsernameCheck model.User
	var placeholderEmailCheck model.User
	lc.DB.Where("username = ?", credentials.Username).First(&placeholderUsernameCheck)
	lc.DB.Where("email = ?", credentials.Username).First(&placeholderEmailCheck)

	if placeholderUsernameCheck.Username != "" {
		w.WriteHeader(http.StatusForbidden)
		log.Printf("Username `%s` already taken", placeholderUsernameCheck.Username)
		return
	}

	if placeholderEmailCheck.Username != "" {
		w.WriteHeader(http.StatusForbidden)
		log.Printf("Email Address `%s` already in use", placeholderEmailCheck.Email)
		return
	}

	if success, _ := regexp.Match("\\s*", []byte(credentials.Password)); !success {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Password too short")
		return
	}

	if !model.ValidateEmail(credentials.Email) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Invalid Email %s", credentials.Email)
		return
	}
	hashedBytes := sha256.Sum256([]byte(credentials.Password + app.GetServerConfig().ServerSettings.Salt))
	hashedPass := string(hashedBytes[:])

	lc.DB.Create(&model.User{
		Username: credentials.Username,
		Email:    credentials.Email,
		Password: hashedPass,
	})

	sessionToken, _ := uuid.NewV4()

	err = app.GetCache().Set(sessionToken.String(), credentials.Username, time.Duration(300*time.Second)).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Path:    "/",
		Expires: time.Now().Add(300 * time.Second),
	})

}

func (lc LoginController) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var updatedUserData loginData

	c, err := r.Cookie("session_token")

	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := c.Value

	err = json.NewDecoder(r.Body).Decode(&updatedUserData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := app.GetCache().Get(sessionToken).Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !model.ValidateEmail(updatedUserData.Email) && updatedUserData.Email != "" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Invalid Email %s", updatedUserData.Email)
		return
	}

	if success, _ := regexp.Match("\\s*", []byte(updatedUserData.Password)); !success && updatedUserData.Password != "" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Password too short")
		return
	}

	var userData model.User

	lc.DB.Where("username = ?", response).First(&userData)

	if updatedUserData.Password != "" {
		newHashedPass := sha256.Sum256([]byte(updatedUserData.Password + app.GetServerConfig().ServerSettings.Salt))
		updatedUserData.Password = string(newHashedPass[:])
	}

	if updatedUserData.Username != "" {
		userData.Username = updatedUserData.Username
	}
	if updatedUserData.Email != "" {
		userData.Email = updatedUserData.Email
	}
	if updatedUserData.Password != "" {
		userData.Password = updatedUserData.Password
	}

	lc.DB.Save(&userData)

	newSessionToken, _ := uuid.NewV4()

	err = app.GetCache().Set(newSessionToken.String(), userData.Username, time.Duration(300*time.Second)).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = app.GetCache().Del(sessionToken).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken.String(),
		Path:    "/",
		Expires: time.Now().Add(300 * time.Second),
	})

}
