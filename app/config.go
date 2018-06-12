package app

import (
	"os"
)

type DatabaseSettings struct {
	User         string
	Password     string
	Address      string
	DatabaseName string
}

type ServerSettings struct {
	Port int
	Salt string
}

type ServerConfig struct {
	ServerSettings   ServerSettings
	DatabaseSettings DatabaseSettings
}

func GenerateConfig() {

	configFile, err := os.Create("./config/config.toml")

	if err != nil {
		panic(err)
	}

	defer configFile.Close()

	initialConfig := `[ServerSettings]
Port = 3000
Salt = ""

[DatabaseSettings]
User = ""
Password = ""
Address = "127.0.0.1"
DatabaseName = "scouting"`
	configFile.Write([]byte(initialConfig))

}
