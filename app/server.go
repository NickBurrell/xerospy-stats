package app

import (
	"context"
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/fatih/color"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/pelletier/go-toml"
	"github.com/zero-frost/xerospy-stats/app/routes"
	"github.com/zero-frost/xerospy-stats/app/routes/middleware"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

type Server struct {
	ServerConfig ServerConfig
	DebugMode    bool
	Cache        *redis.Client
	Router       *mux.Router
	Database     *gorm.DB
}

func (s *Server) InitServer() {
	if _, err := os.Stat("config/config.toml"); os.IsNotExist(err) {
		if s.DebugMode {
			fmt.Print(color.YellowString("[!] Configuration file not found. Generating new configuration file\n"))
		}
		GenerateConfig()
	}
	if s.DebugMode {
		fmt.Print(color.GreenString("[*] Configuration file found\n"))
	}
	file, err := os.Open("./config/config.toml")
	if err != nil {
		panic(err)
	}
	fileInfo, _ := file.Stat()

	buffer := make([]byte, fileInfo.Size())
	file.Read(buffer)

	toml.Unmarshal(buffer, &s.ServerConfig)

	// Initialize the templating engine and data storage systems
	s.InitHandlebars()
	s.InitDatabase()
	s.InitRedis()

	if _, err := os.Stat("./log"); os.IsNotExist(err) {
		os.Mkdir("./log", 0755)
	}

	if _, err = os.Stat("./log/" + s.ServerConfig.ServerSettings.LogFile); err != nil {
		os.Create("./log/" + s.ServerConfig.ServerSettings.LogFile)
	}

	logFile, err := os.OpenFile("log/"+s.ServerConfig.ServerSettings.LogFile, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}

	if s.DebugMode {
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
	} else {
		log.SetOutput(logFile)
	}

	if s.ServerConfig.ServerSettings.Salt == "" {
		fmt.Print(color.RedString("[!] Error: Salt not set! Please set a new password salt to continue\n"))
		os.Exit(3)
	}

	s.Router = mux.NewRouter()

	apiHandler := routes.GetAPIRoutes(s.Database, s.Cache, s.ServerConfig.ServerSettings.Salt)

	s.Router.Handle("/api/{_dummy:.*}", apiHandler)

	// Routes for static assets
	s.Router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	s.Router.PathPrefix("/lib/").Handler(http.StripPrefix("/lib/", http.FileServer(http.Dir("node_modules"))))
	s.Router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Define middlewares
	s.Router.Use(middleware.LoggingMiddleware)

	log.Println("Server Initialized")

}

func (s *Server) RunServer() {

	srv := &http.Server{
		Addr:         "0.0.0.0:" + strconv.Itoa(s.ServerConfig.ServerSettings.Port),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.Router,
	}

	// Serve HTTP
	if s.DebugMode {
		fmt.Printf(color.GreenString("[*] Listening on port :%d\n", s.ServerConfig.ServerSettings.Port))
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
		defer s.Database.Close()
	}()

	log.Println("Server Started")

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	srv.Shutdown(ctx)
	if s.DebugMode {
		fmt.Println(color.YellowString("[!] Server Shutting Down\n"))
	}
	log.Println("Server Shutting Down")
	os.Exit(0)

}

func (s *Server) InitHandlebars() error {

	header, err := ioutil.ReadFile("template/header.hbs")
	if err != nil {
		return err
	}

	raymond.RegisterPartial("header", string(header))

	if s.DebugMode {
		fmt.Print(color.GreenString("[*] Registered Partial: header\n"))
	}

	return nil
}

func (s *Server) RenderTemplate(filename string, params map[string]string) (string, error) {
	file, err := ioutil.ReadFile("template/" + filename)
	if err != nil {
		return "", err
	}
	contents, err := raymond.Render(string(file), params)
	if err != nil {
		return "", err
	}
	return contents, nil
}

func (s *Server) InitRedis() {

	s.Cache = redis.NewClient(&redis.Options{
		Addr:     s.ServerConfig.RedisSettings.Address,
		Password: s.ServerConfig.RedisSettings.Password,
		DB:       s.ServerConfig.RedisSettings.Database,
	})
}
