package controllers

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/zero-frost/xerospy-stats/app/model"
	"log"
	"net/http"
	"regexp"
	"time"
)

type LoginController struct {
	*Controller
	Salt string
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

	hashedPass := sha256.Sum256([]byte(credentials.Password + lc.Salt))

	if !bytes.Equal([]byte(userData.Password), hashedPass[:]) {
		fmt.Fprint(w, "Incorrect Password\n")
		return
	}

	sessionToken, _ := uuid.NewV4()

	err = lc.Cache.Set(sessionToken.String(), credentials.Username, time.Duration(300*time.Second)).Err()
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

func (lc LoginController) Refresh(w http.ResponseWriter, r *http.Request) {

	var sessionToken string
	var err error
	var statusCode int
	if statusCode, sessionToken, err = lc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	res, _ := lc.Cache.Get(sessionToken).Result()

	newSessionToken, _ := uuid.NewV4()
	err = lc.Cache.Set(newSessionToken.String(), fmt.Sprint(res), time.Duration(300*time.Second)).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = lc.Cache.Del(sessionToken).Err()
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

	var sessionToken string
	var err error
	var statusCode int
	if statusCode, sessionToken, err = lc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	err = lc.Cache.Del(sessionToken).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}

func (lc LoginController) CreateUser(w http.ResponseWriter, r *http.Request) {

	var credentials loginData

	var err error
	if err = json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if credentials.Email == "" || credentials.Password == "" || credentials.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
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

	var placeholderUsernameCheck model.User
	lc.DB.Where("username = ?", credentials.Username).First(&placeholderUsernameCheck)

	if placeholderUsernameCheck.Username != "" {
		w.WriteHeader(http.StatusForbidden)
		log.Printf("Username `%s` already taken", placeholderUsernameCheck.Username)
		return
	}

	var placeholderEmailCheck model.User
	lc.DB.Where("email = ?", credentials.Username).First(&placeholderEmailCheck)

	if placeholderEmailCheck.Username != "" {
		w.WriteHeader(http.StatusForbidden)
		log.Printf("Email Address `%s` already in use", placeholderEmailCheck.Email)
		return
	}

	hashedBytes := sha256.Sum256([]byte(credentials.Password + lc.Salt))
	hashedPass := string(hashedBytes[:])

	lc.DB.Create(&model.User{
		Username: credentials.Username,
		Email:    credentials.Email,
		Password: hashedPass,
	})

}

func (lc LoginController) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var sessionToken string
	var err error
	var statusCode int
	if statusCode, sessionToken, err = lc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	var updatedUserData loginData

	err = json.NewDecoder(r.Body).Decode(&updatedUserData)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := lc.Cache.Get(sessionToken).Result()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if success, _ := regexp.Match("\\s*", []byte(updatedUserData.Password)); !success && updatedUserData.Password != "" {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Password too short")
		return
	}

	var userData model.User

	lc.DB.Where("username = ?", response).First(&userData)

	err = json.NewDecoder(r.Body).Decode(&userData)

	if updatedUserData.Password != "" {
		newHashedPass := sha256.Sum256([]byte(updatedUserData.Password + lc.Salt))
		updatedUserData.Password = string(newHashedPass[:])
	}

	lc.DB.Save(&userData)

	newSessionToken, _ := uuid.NewV4()

	err = lc.Cache.Set(newSessionToken.String(), userData.Username, time.Duration(300*time.Second)).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = lc.Cache.Del(sessionToken).Err()
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

func (lc LoginController) DeleteUser(w http.ResponseWriter, r *http.Request) {

	var sessionToken string
	var err error
	var statusCode int
	if statusCode, sessionToken, err = lc.ValidateSession(r); err != nil {
		w.WriteHeader(statusCode)
		return
	}

	var credentials loginData

	err = json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = lc.Cache.Get(sessionToken).Err()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var userData model.User
	lc.DB.Where("username = ?", credentials.Username).First(&userData)

	if userData.Username == "" {
		fmt.Fprintf(w, "No user found by the username %s\n", credentials.Username)
		return
	}

	hashedPass := sha256.Sum256([]byte(credentials.Password + lc.Salt))

	if !bytes.Equal([]byte(userData.Password), hashedPass[:]) {
		fmt.Fprint(w, "Incorrect Password\n")
		return
	}

	err = lc.DB.Delete(userData).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Path:    "/",
		MaxAge:  -1,
		Expires: time.Now(),
	})
}
