package main

import (
	"log"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

func main() {

	e := echo.New()
	e.Debug()

	//Midleware
	e.Use(mw.Logger())
	e.Use(mw.Recover())

	e.SetHTTPErrorHandler(func(err error, c *echo.Context) {
		log.Println(err)
		e.DefaultHTTPErrorHandler(err, c)
	})

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
