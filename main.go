package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/flexicon/bookscale/views"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
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

	if err := InitCache(); err != nil {
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
	SetupRoutes(e)

	return e.Start(fmt.Sprintf(":%d", viper.GetInt("port")))
}

// ViperInit loads environment variables and sets up needed defaults.
func ViperInit() error {
	// Prepare for Environment variables
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Defaults
	viper.SetDefault("port", 80)
	viper.SetDefault("debug", false)
	viper.SetDefault("cache.ttl", 900) // In seconds
	viper.SetDefault("static_asset_base_url", "https://res.cloudinary.com/flexicondev/image/upload/v1646583872/bookscale/")
	viper.SetDefault("allegro.client_id", "")
	viper.SetDefault("allegro.client_secret", "")

	// Read optional config.yml file
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return errors.Wrap(err, "failed to read existing viper config file")
		}
	}

	return nil
}
