package main

import (
	"flag"
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm"
	"github.com/pelletier/go-toml"
	"github.com/zero-frost/xerospy-stats/app"
	_ "github.com/zero-frost/xerospy-stats/app/model"
	"github.com/zero-frost/xerospy-stats/app/routes"
	"github.com/zero-frost/xerospy-stats/app/routes/middleware"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var isDebug bool = false

func initServer() app.ServerConfig {
	if _, err := os.Stat("config/config.toml"); os.IsNotExist(err) {
		if isDebug {
			fmt.Print(color.YellowString("[!] Configuration file not found. Generating new configuration file\n"))
		}
		app.GenerateConfig()
	}
	if isDebug {
		fmt.Print(color.GreenString("[*] Configuration file found\n"))
	}
	file, err := os.Open("./config/config.toml")
	if err != nil {
		panic(err)
	}
	fileInfo, _ := file.Stat()

	buffer := make([]byte, fileInfo.Size())
	file.Read(buffer)

	config := app.ServerConfig{}
	toml.Unmarshal(buffer, &config)

	if config.ServerSettings.Salt == "" {
		fmt.Print(color.RedString("[!] Error: Salt not set! Please set a new password salt to continue\n"))
		os.Exit(3)
	}

	return config
}

func initHandlebars() error {

	header, err := ioutil.ReadFile("template/header.hbs")
	if err != nil {
		return err
	}

	raymond.RegisterPartial("header", string(header))

	if isDebug {
		fmt.Print(color.GreenString("[*] Registered Partial: header\n"))
	}

	return nil
}

func main() {

	// Parse command line flags
	flag.BoolVar(&isDebug, "debug", false, "Enables debug print-outs")
	flag.Parse()

	// Initialize server and templating engine
	config := initServer()
	initHandlebars()
	app.InitDatabase(config.DatabaseSettings.User,
		config.DatabaseSettings.Password,
		config.DatabaseSettings.Address,
		config.DatabaseSettings.DatabaseName)

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
	if isDebug {
		fmt.Printf(color.GreenString("[*] Listening on port :%d\n", config.ServerSettings.Port))
	}

	http.ListenAndServe(":"+strconv.Itoa(config.ServerSettings.Port), r)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	file, err := ioutil.ReadFile("template/index.hbs")
	if err != nil {
		panic(err)
	}
	index, err := raymond.Render(string(file), map[string]string{"title": "Hello"})
	if err != nil {
		panic(err)
	}
	fmt.Fprint(w, index)

}
