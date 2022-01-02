package main

import (
	"log"
	"strings"

	"github.com/flexicon/bookscale/views"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	if err := ViperInit(); err != nil {
		return err
	}

	e := echo.New()
	e.Debug = viper.GetBool("debug")

	e.Use(middleware.Secure())
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "REQUEST: method=${method}, status=${status}, uri=${uri}, latency=${latency_human}\n",
	}))

	e.Renderer = views.NewRenderer()

	e.GET("/", IndexHandler)
	e.GET("/search", SearchHandler)

	return e.Start(":9000")
}

// ViperInit loads environment variables and sets up needed defaults.
func ViperInit() error {
	// Prepare for Environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("port", 80)
	viper.SetDefault("debug", false)
	viper.SetDefault("allegro.client_id", "")
	viper.SetDefault("allegro.client_secret", "")

	return nil
}
