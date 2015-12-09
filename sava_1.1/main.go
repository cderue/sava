package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"crypto/rand"
	"fmt"
	"syscall"
	"os/signal"
)

var (
	uuid = generateUuid()
	debug = flag.Bool("debug", false, "Sets the debug mode")
	portHtml = flag.Int("html", 8080, "Sets the port to listen on")
	portJson = flag.Int("json", 8081, "Sets the port to listen on")
)

func main() {

	parseParameters()

	go serverHtml()
	go serverJson()

	// Waiter keeps the program from exiting instantly.
	waiter := make(chan bool)

	// Catch a CTR+C exits so the cleanup routine is called.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(1)
	}()

	waiter <- true
}

func parseParameters() {

	flag.Parse()

	envPortHtml := os.Getenv("SAVA_PORT_HTML")

	if envPortHtml != "" {
		*portHtml, _ = strconv.Atoi(envPortHtml)
	}

	envPortJson := os.Getenv("SAVA_PORT_JSON")

	if envPortJson != "" {
		*portJson, _ = strconv.Atoi(envPortJson)
	}

	envDebug := os.Getenv("SAVA_DEBUG")

	if envDebug != "" {
		*debug, _ = strconv.ParseBool(envDebug)
	}
}

func serverHtml() {

	router := gin.Default()

	router.LoadHTMLGlob("index.tmpl")
	router.Static("/public/js", "./public/js")
	router.Static("/public/css", "./public/css")

	router.GET("/", func(c *gin.Context) {
		var display = "none"
		if *debug {
			display = "block"
		}
		obj := gin.H{"display": display, "version": version, "uuid": uuid}
		c.HTML(200, "index.tmpl", obj)
	})

	router.GET("/public", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/")
	})

	router.GET("public/favicon.png", func(c *gin.Context) {
		c.File("./public/favicon.png")
	})

	router.Run(":" + strconv.Itoa(*portHtml))
}

func serverJson() {

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": version, "uuid": uuid})
	})

	router.Run(":" + strconv.Itoa(*portJson))
}


func generateUuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}
