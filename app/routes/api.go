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

	mc := controllers.MatchController{
		Controller: &controllers.Controller{
			DB:    db,
			Cache: cache,
		},
	}

	atc := controllers.ActionTypeController{
		Controller: &controllers.Controller{
			DB:    db,
			Cache: cache,
		},
	}

	tmc := controllers.TeamMatchController{
		Controller: &controllers.Controller{
			DB:    db,
			Cache: cache,
		},
	}

	tmac := controllers.TeamMatchActionController{
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

	// Regiser EventController routes
	apiRouter.HandleFunc("/api/event/{event_id}", ec.GetEvent).Methods("GET")
	apiRouter.HandleFunc("/api/event", ec.GetEvent).Methods("GET")
	apiRouter.HandleFunc("/api/event/{event_id}", ec.UpdateEvent).Methods("PUT")
	apiRouter.HandleFunc("/api/event", ec.AddEvent).Methods("POST")
	apiRouter.HandleFunc("/api/event", ec.DeleteEvent).Methods("DELETE")

	// Regiser MatchController routes
	apiRouter.HandleFunc("/api/match/{match_id}", mc.GetMatch).Methods("GET")
	apiRouter.HandleFunc("/api/match", mc.GetMatches).Methods("GET")
	apiRouter.HandleFunc("/api/match/{event_id}", mc.UpdateMatch).Methods("PUT")
	apiRouter.HandleFunc("/api/match", mc.AddMatch).Methods("POST")
	apiRouter.HandleFunc("/api/match", mc.DeleteMatch).Methods("DELETE")

	// Register EventController routes
	apiRouter.HandleFunc("/api/team/{team_id}", tc.GetTeam).Methods("GET")
	apiRouter.HandleFunc("/api/team", tc.GetTeams).Methods("GET")
	apiRouter.HandleFunc("/api/team/{team_id}", tc.UpdateTeam).Methods("PUT")
	apiRouter.HandleFunc("/api/team", tc.UpdateTeam).Methods("POST")
	apiRouter.HandleFunc("/api/team", tc.DeleteTeam).Methods("DELETE")

	// Regiser ActionTypeController routes
	apiRouter.HandleFunc("/api/action_type/{action_type_id}", atc.GetActionType).Methods("GET")
	apiRouter.HandleFunc("/api/action_type", atc.GetActionTypes).Methods("GET")
	apiRouter.HandleFunc("/api/action_type/{action_type_id}", atc.UpdateActionType).Methods("PUT")
	apiRouter.HandleFunc("/api/action_type", atc.AddActionType).Methods("POST")
	apiRouter.HandleFunc("/api/action_type", atc.DeleteActionType).Methods("DELETE")

	// Regiser TeamMatchController routes
	apiRouter.HandleFunc("/api/team_match/{team_match_id}", tmc.GetTeamMatch).Methods("GET")
	apiRouter.HandleFunc("/api/team_match", tmc.GetTeamMatches).Methods("GET")
	apiRouter.HandleFunc("/api/team_match/{team_match_id}", tmc.UpdateTeamMatch).Methods("PUT")
	apiRouter.HandleFunc("/api/team_match", tmc.AddTeamMatch).Methods("POST")
	apiRouter.HandleFunc("/api/team_match", tmc.DeleteTeamMatch).Methods("DELETE")

	// Regiser TeamMatchActionController routes
	apiRouter.HandleFunc("/api/team_match_action/{team_match_action_id}", tmac.GetTeamMatchAction).Methods("GET")
	apiRouter.HandleFunc("/api/team_match_action", tmac.GetTeamMatchAction).Methods("GET")
	apiRouter.HandleFunc("/api/team_match_action/{team_match_action_id}", tmac.UpdateTeamMatchAction).Methods("PUT")
	apiRouter.HandleFunc("/api/team_match_action", tmac.AddTeamMatchAction).Methods("POST")
	apiRouter.HandleFunc("/api/team_match_action", tmac.DeleteTeamMatchAction).Methods("DELETE")

	// Register LoginController routes
	apiRouter.HandleFunc("/api/user/login", lc.Login).Methods("POST")
	apiRouter.HandleFunc("/api/user/refresh", lc.Refresh)
	apiRouter.HandleFunc("/api/user/logout", lc.Logout)
	apiRouter.HandleFunc("/api/user", lc.CreateUser).Methods("POST")
	apiRouter.HandleFunc("/api/user", lc.UpdateUser).Methods("PUT")
	apiRouter.HandleFunc("/api/user", lc.DeleteUser)

	return apiRouter
}
