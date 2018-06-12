package routes

import (
	"github.com/gorilla/mux"
	"github.com/zero-frost/xerospy-stats/app"
	"github.com/zero-frost/xerospy-stats/app/controllers/api"
)

func GetAPIRoutes() *mux.Router {

	tc := api.TeamController{DB: app.GetDatabase()}

	apiRouter := mux.NewRouter()

	apiRouter.HandleFunc("/api/team/{team_id}", tc.GetTeam).Methods("GET")
	apiRouter.HandleFunc("/api/team", tc.GetTeams).Methods("GET")
	return apiRouter
}
