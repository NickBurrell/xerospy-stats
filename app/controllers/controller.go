package controllers

import (
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"net/http"
)

type Controller struct {
	DB    *gorm.DB
	Cache *redis.Client
}

func (c *Controller) ValidateSession(r *http.Request) (int, string, error) {

	cookie, err := r.Cookie("session_token")

	if err != nil {
		if err == http.ErrNoCookie {
			return http.StatusUnauthorized, "", http.ErrNoCookie
		}
		return http.StatusUnauthorized, "", http.ErrNoCookie
	}

	if err != nil {
		return http.StatusUnauthorized, "", err
	}

	sessionToken := cookie.Value

	var res string
	res, err = c.Cache.Get(sessionToken).Result()

	if err != nil {
		return http.StatusInternalServerError, "", err
	}

	if res == "" {
		return http.StatusUnauthorized, "", err
	}

	return http.StatusOK, sessionToken, nil

}
