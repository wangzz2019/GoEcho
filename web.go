package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", hello)
	e.GET("/test", test)

	e.GET("/user/:id", getUser)

	e.File("/home", "goweb.html")
	e.Static("/img", "img")
	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World! This is jack's first echo web!")
}
func test(c echo.Context) error {
	return c.String(http.StatusOK, "Hi, This is a test page")
}

func getUser(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, "User Id is: "+id)
}
