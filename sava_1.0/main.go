package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

var (
	port = flag.Int("port", 8080, "Sets the port to listen on")
	debug = flag.Bool("debug", false, "Sets the debug mode")
)

func main() {
	parameters()
	server()
}

func parameters() {

	flag.Parse()

	envPort := os.Getenv("SAVA_PORT")

	if envPort != "" {
		*port, _ = strconv.Atoi(envPort)
	}

	envDebug := os.Getenv("SAVA_DEBUG")

	if envDebug != "" {
		*debug, _ = strconv.ParseBool(envDebug)
	}
}

func server() {

	r := gin.Default()
	r.Static("/public", "./public")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/public")
	})

	r.Run(":" + strconv.Itoa(*port))
}
