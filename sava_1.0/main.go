package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
	"crypto/rand"
	"fmt"
)

var (
	version = "sava:1.0"
	uuid = generate_uuid()
	debug = flag.Bool("debug", false, "Sets the debug mode")
	port = flag.Int("port", 8080, "Sets the port to listen on")
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

	router.Run(":" + strconv.Itoa(*port))
}

func generate_uuid() (uuid string) {

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	uuid = fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
	return
}
