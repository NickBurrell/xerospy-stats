package routes

import (
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/zero-frost/xerospy-stats/app/controllers"
)

func GetAPIRoutes(db *gorm.DB, cache *redis.Client, salt string) *mux.Router {

	tc := controllers.TeamController{
		Controller: &controllers.Controller{
			DB:    db,
			Cache: cache,
		},
	}
	ec := controllers.EventController{
		Controller: &controllers.Controller{
			DB:    db,
			Cache: cache,
		},
	}
	lc := controllers.LoginController{
		Controller: &controllers.Controller{
			DB:    db,
			Cache: cache,
		},
		Salt: salt,
	}

	apiRouter := mux.NewRouter()

	// Regiser TeamController routes
	apiRouter.HandleFunc("/api/team/{team_id}", ec.GetEvent).Methods("GET")
	apiRouter.HandleFunc("/api/team", ec.GetEvent).Methods("GET")
	apiRouter.HandleFunc("/api/team/{team_id}", ec.UpdateEvent).Methods("PUT")
	apiRouter.HandleFunc("/api/team", ec.UpdateEvent).Methods("POST")
	apiRouter.HandleFunc("/api/team", ec.DeleteEvent).Methods("DELETE")

	// Register EventController routes
	apiRouter.HandleFunc("/api/event/{team_id}", tc.GetTeam).Methods("GET")
	apiRouter.HandleFunc("/api/event", tc.GetTeams).Methods("GET")
	apiRouter.HandleFunc("/api/event/{team_id}", tc.UpdateTeam).Methods("PUT")
	apiRouter.HandleFunc("/api/event", tc.UpdateTeam).Methods("POST")
	apiRouter.HandleFunc("/api/event", tc.DeleteTeam).Methods("DELETE")

	// Register LoginController routes
	apiRouter.HandleFunc("/api/user/login", lc.Login).Methods("POST")
	apiRouter.HandleFunc("/api/user/refresh", lc.Refresh)
	apiRouter.HandleFunc("/api/user/logout", lc.Logout)
	apiRouter.HandleFunc("/api/user", lc.CreateUser).Methods("POST")
	apiRouter.HandleFunc("/api/user", lc.UpdateUser).Methods("PUT")
	apiRouter.HandleFunc("/api/user", lc.DeleteUser)

	return apiRouter
}
