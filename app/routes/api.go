package routes

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/zero-frost/xerospy-stats/app/controllers"
	"net/http"
)

func GetAPIRoutes() *mux.Router {

	apiRouter := mux.NewRouter()

	apiRouter.HandleFunc("/api/team/{team_id}", controllers.GetTeam).Methods("GET")
	apiRouter.HandleFunc("/api/team", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Test")
	})
	return apiRouter
}
