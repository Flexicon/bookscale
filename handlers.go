package main

import (
	"net/http"

	"github.com/flexicon/bookscale/views"
	"github.com/labstack/echo/v4"
)

func IndexHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", nil)
}

func SearchHandler(c echo.Context) error {
	searchTerm := c.QueryParam("term")

	return c.Render(http.StatusOK, "search.html", views.Args{
		"SearchTerm": searchTerm,
	})
}
