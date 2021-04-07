package main

import (
	"net/http"
	"os"
	"fmt"
	"time"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	datadog "github.com/DataDog/datadog-api-client-go/api/v2/datadog"
	"context"
)

type AccountType struct {
	ID   string `json:"id" form:"id" query:"id"`
	PASS string `json:"password" form:"password" query:"password"`
}

func main() {
	// use JSONFormatter
	log.SetFormatter(&log.JSONFormatter{})
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Info("Failed to log to file")
	}
	tracer.Start(
		tracer.WithEnv("prod"),
		tracer.WithService("goweb"),
		tracer.WithDebugMode(true),
		//tracer.WithVersion("abc123"),
	)

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

	//Post
	e.POST("/gettoken", gettoken)
	e.POST("/checktoken", checktoken)
	e.POST("/webhook",webhook)
	e.POST("/ddapi",ddapi)
	e.Static("/img", "img")
	// Start server
	e.Logger.Fatal(e.Start(":1323"))
	defer tracer.Stop()
}

// Handler
func webhook(c echo.Context) (err error){
	return nil
}
func ddapi(c echo.Context) (err error){
	ctx := datadog.NewDefaultContext(context.Background())
	filterQuery := "filename:jack.log" // string | Search query following logs syntax. (optional)
    filterIndex := "main" // string | For customers with multiple indexes, the indexes to search Defaults to '*' which means all indexes (optional)
    filterFrom := time.Now().Add(-time.Minute * 15) // time.Time | Minimum timestamp for requested logs. (optional)
    filterTo := time.Now() // time.Time | Maximum timestamp for requested logs. (optional)
    sort := datadog.LogsSort("timestamp") // LogsSort | Order of logs in results. (optional)
    // pageCursor := "eyJzdGFydEF0IjoiQVFBQUFYS2tMS3pPbm40NGV3QUFBQUJCV0V0clRFdDZVbG8zY3pCRmNsbHJiVmxDWlEifQ==" // string | List following results with a cursor provided in the previous query. (optional)
    pageLimit := int32(25) // int32 | Maximum number of logs in the response. (optional) (default to 10)

    configuration := datadog.NewConfiguration()

    apiClient := datadog.NewAPIClient(configuration)
    resp, r, err := apiClient.LogsApi.ListLogsGet(ctx).FilterQuery(filterQuery).FilterIndex(filterIndex).FilterFrom(filterFrom).FilterTo(filterTo).Sort(sort).PageLimit(pageLimit).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `LogsApi.ListLogsGet``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ListLogsGet`: LogsListResponse
    responseContent, _ := json.MarshalIndent(resp, "", "  ")
    fmt.Fprintf(os.Stdout, "Response from LogsApi.ListLogsGet:\n%s\n", responseContent)

	return nil
}

func gettoken(c echo.Context) (err error) {
	Account := new(AccountType)
	if err = c.Bind(Account); err != nil {
		return
	}
	var Token string
	if Account.ID == "jack" && Account.PASS == "wang" {
		Token = "TokenStringABCDEFG"
		// return c.JSON(http.StatusOK, Token)
		return c.String(http.StatusOK, Token)
	}
	Token = "error!"
	return c.JSON(http.StatusNonAuthoritativeInfo, Token)
	// return c.JSON(http.StatusOK, Token)
}
func checktoken(c echo.Context) error {
	token := c.Request().Header.Get("Authorization")
	if token == "TokenStringABCDEFG" {
		return c.JSON(http.StatusOK, "PASS!")
	}
	return c.JSON(http.StatusUnauthorized, "Not Pass!")
}
func hello(c echo.Context) error {
	// log an event as usual with logrus
	log.WithFields(log.Fields{"string": "foo", "int": 1, "float": 1.1}).Info("My first event from golang to stdout")
	// For metadata, a common pattern is to re-use fields between logging statements  by re-using
	contextualizedLog := log.WithFields(log.Fields{
		"hostname": "staging-1",
		"appname":  "foo-app",
		"session":  "1ce3f6v"})

	contextualizedLog.Info("Simple event with global metadata")
	return c.String(http.StatusOK, "Hello, World! This is jack's first echo web!")
}
func test(c echo.Context) error {
	return c.String(http.StatusOK, "Hi, This is a test page")
}

func getUser(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, "User Id is: "+id)
}
