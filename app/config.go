package app

import (
	"os"
)

type DatabaseSettings struct {
	Address      string
	Username     string
	Password     string
	DatabaseName string
}

type ServerSettings struct {
	Port int
	Salt string
}

type RedisSettings struct {
	Address  string
	Password string
	Database int
}

type ServerConfig struct {
	ServerSettings   ServerSettings
	DatabaseSettings DatabaseSettings
	RedisSettings    RedisSettings
}

func GenerateConfig() {

	if _, err := os.Stat("./config"); os.IsNotExist(err) {
		os.Mkdir("./config", 0755)
	}

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
Address = "127.0.0.1:3306"
DatabaseName = "scouting"

[RedisSettings]
Address = "127.0.0.1:6379"
Password = ""
Database = 0
`
	configFile.Write([]byte(initialConfig))

}
