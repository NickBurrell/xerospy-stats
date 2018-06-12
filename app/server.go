package app

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/pelletier/go-toml"
	"os"
)

var DebugMode bool = false

var serverConfig ServerConfig

func InitServer() {
	if _, err := os.Stat("config/config.toml"); os.IsNotExist(err) {
		if DebugMode {
			fmt.Print(color.YellowString("[!] Configuration file not found. Generating new configuration file\n"))
		}
		GenerateConfig()
	}
	if DebugMode {
		fmt.Print(color.GreenString("[*] Configuration file found\n"))
	}
	file, err := os.Open("./config/config.toml")
	if err != nil {
		panic(err)
	}
	fileInfo, _ := file.Stat()

	buffer := make([]byte, fileInfo.Size())
	file.Read(buffer)

	config := ServerConfig{}
	toml.Unmarshal(buffer, &config)

	serverConfig = config
}

func GetServerConfig() *ServerConfig {
	return &serverConfig
}
