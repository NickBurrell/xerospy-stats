package app

import (
	"fmt"
	"github.com/aymerick/raymond"
	"github.com/fatih/color"
	"io/ioutil"
)

func InitHandlebars() error {

	header, err := ioutil.ReadFile("template/header.hbs")
	if err != nil {
		return err
	}

	raymond.RegisterPartial("header", string(header))

	if DebugMode {
		fmt.Print(color.GreenString("[*] Registered Partial: header\n"))
	}

	return nil
}

func RenderTemplate(filename string, params map[string]string) (string, error) {
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
