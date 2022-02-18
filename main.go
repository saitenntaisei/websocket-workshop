package main

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()

	e.Logger.SetLevel(log.DEBUG)
	e.Logger.SetHeader("${time_rfc3339} ${prefix} ${short_file} ${line} |")
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: "${time_rfc3339} method = ${method} | uri = ${uri} | code = ${status} ${error}\n"}))

	e.Static("/", "public")

	api := e.Group("/api")
	{
		api.GET("/ping", func(c echo.Context) error {
			return c.String(http.StatusOK, "pong")
		})
		api.GET("/ws", connectWS)
	}

	e.Logger.Panic(e.Start(":8080"))
}

var upgrader = websocket.Upgrader{}

func connectWS(c echo.Context) error {
	connection, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	defer connection.Close()
	for {
		mesType, msg, err := connection.ReadMessage()
		if err != nil {
			return c.NoContent(http.StatusOK)
		}
		log.Printf("mesType= %d,msg=%s,", mesType, msg)
		err = connection.WriteMessage(websocket.TextMessage, []byte("Hello"))
		if err != nil {
			return c.NoContent((http.StatusOK))
		}
	}
}
