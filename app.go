package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm"
	"github.com/zero-frost/xerospy-stats/app"
	_ "github.com/zero-frost/xerospy-stats/app/model"
	"github.com/zero-frost/xerospy-stats/app/routes"
	"github.com/zero-frost/xerospy-stats/app/routes/middleware"
	"net/http"
	"os"
	"strconv"
)

func main() {

	// Parse command line flags
	flag.BoolVar(&app.DebugMode, "debug", false, "Enables debug print-outs")
	flag.Parse()

	// Initialize server and templating engine
	app.InitServer()
	config := app.GetServerConfig()

	// Initialize the templating engine
	app.InitHandlebars()

	app.InitDatabase(config.DatabaseSettings.Username,
		config.DatabaseSettings.Password,
		config.DatabaseSettings.Address,
		config.DatabaseSettings.DatabaseName)

	app.InitRedis(config.RedisSettings.Address,
		config.RedisSettings.Password,
		config.RedisSettings.Database)

	if app.GetServerConfig().ServerSettings.Salt == "" {
		fmt.Print(color.RedString("[!] Error: Salt not set! Please set a new password salt to continue\n"))
		os.Exit(3)
	}

	defer app.GetDatabase().Close()

	r := mux.NewRouter()

	// Define middlewares
	r.Use(middleware.LoggingMiddleware)

	// Set up Routes
	r.HandleFunc("/", handleRoot)

	apiRouter := routes.GetAPIRoutes()

	r.Handle("/api/{_dummy:.*}", apiRouter)

	// Routes for static assets
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.PathPrefix("/lib/").Handler(http.StripPrefix("/lib/", http.FileServer(http.Dir("node_modules"))))
	r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Serve HTTP
	if app.DebugMode {
		fmt.Printf(color.GreenString("[*] Listening on port :%d\n", config.ServerSettings.Port))
	}

	http.ListenAndServe(":"+strconv.Itoa(config.ServerSettings.Port), r)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	index, err := app.RenderTemplate("index.hbs",
		map[string]string{
			"title": "Hello!",
		})
	if err != nil {
		fmt.Fprint(w, http.StatusInternalServerError)
	}
	fmt.Fprint(w, index)

}
