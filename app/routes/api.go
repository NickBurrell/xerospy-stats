package routes

import (
	"github.com/gorilla/mux"
	"github.com/zero-frost/xerospy-stats/app"
	"github.com/zero-frost/xerospy-stats/app/controllers"
	"github.com/zero-frost/xerospy-stats/app/controllers/api"
)

func GetAPIRoutes() *mux.Router {

	tc := api.TeamController{DB: app.GetDatabase()}
	lc := controllers.LoginController{DB: app.GetDatabase()}

	apiRouter := mux.NewRouter()

	apiRouter.HandleFunc("/api/team/{team_id}", tc.GetTeam).Methods("GET")
	apiRouter.HandleFunc("/api/team", tc.GetTeams).Methods("GET")

	// TEMPORARY
	apiRouter.HandleFunc("/api/user/login", lc.Login).Methods("POST")
	apiRouter.HandleFunc("/api/user/refresh", lc.Refresh)
	apiRouter.HandleFunc("/api/user/logout", lc.Logout)
	return apiRouter
}
