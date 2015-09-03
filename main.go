package main

import (
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

func main() {
	e := echo.New()

	//Midleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	//Handlers
	e.Post("/Plugin.Activate", PluginActivate)
	e.Post("/NetworkDriver.CreateNetwork", NetworkDriverCreateNetwork)
	e.Post("/NetworkDriver.DeleteNetwork", NetworkDriverDeleteNetwork)
	e.Post("/NetworkDriver.CreateEndpoint", NetworkDriverCreateEndpoint)
	e.Post("/NetworkDriver.EndpointOperInfo", NetworkDriverEndpointOperInfo)
	e.Post("/NetworkDriver.DeleteEndpoint", NetworkDriverDeleteEndpoint)
	e.Post("/NetworkDriver.Join", NetworkDriverJoin)
	e.Post("/NetworkDriver.Leave", NetworkDriverLeave)

	e.Run(":1313")
}
