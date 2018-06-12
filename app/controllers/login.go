package controllers

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"github.com/zero-frost/xerospy-stats/app"
	"github.com/zero-frost/xerospy-stats/app/model"
	"net/http"
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

	hasher := sha256.New()
	hasher.Write([]byte("D3@ThC0D3error-code-xero"))
	var hexData []byte
	hashedPass := sha256.Sum256([]byte(credentials.Password + app.GetServerConfig().ServerSettings.Salt))

	//hex.Decode(hexData, temp[:])

	if bytes.Equal([]byte(userData.Password), hashedPass[:]) {
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
		Expires: time.Now().Add(300 * time.Second),
	})

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

	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Unix(0, 0),
	})
}
