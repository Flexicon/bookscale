package main

import (
	"log"

	"github.com/flexicon/bookscale/views"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	if err := run(); err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	e := echo.New()
	e.Debug = true // TODO: move this to an env var

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
