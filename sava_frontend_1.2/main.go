package main

import (
	"encoding/json"
	"flag"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var (
	lorem = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
)

func main() {

	port := flag.Int("port", 8080, "Sets the port to listen on")

	flag.Parse()

	envBackend1 := os.Getenv("BACKEND_1")
	envBackend2 := os.Getenv("BACKEND_2")
	envPort := os.Getenv("SAVA_PORT")

	if envPort != "" {
		*port, _ = strconv.Atoi(envPort)
	}

	r := gin.Default()

	r.Static("/js", "./public/js")
	r.Static("/css", "./public/css")
	r.Static("/public/js", "./public/js")
	r.Static("/public/css", "./public/css")

	r.GET("/api/message1", func(c *gin.Context) {
		GetMessage(c, envBackend1)
	})

	r.GET("/public/api/message1", func(c *gin.Context) {
		GetMessage(c, envBackend1)
	})

	r.GET("/api/message2", func(c *gin.Context) {
		GetMessage(c, envBackend2)
	})

	r.GET("/public/api/message2", func(c *gin.Context) {
		GetMessage(c, envBackend2)
	})

	r.GET("/", func(c *gin.Context) {
		c.File("/public/index.html")
	})

	r.GET("/public", func(c *gin.Context) {
		c.File("/public/index.html")
	})

	r.GET("/public/", func(c *gin.Context) {
		c.File("/public/index.html")
	})

	r.Run(":" + strconv.Itoa(*port))
}

func GetMessage(c *gin.Context, backend string) {

	resp, err := http.Get(backend)
	if err != nil {
		c.JSON(200, gin.H{"text": "" + lorem + ""})
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			c.JSON(200, gin.H{"text": "" + lorem + ""})
		} else {
			var data Lorem
			json.Unmarshal(body, &data)

			c.JSON(200, gin.H{"text": "" + data.Text + ""})
		}
	}
}

type Lorem struct {
	Text  string      `json: "text"`
	Paras interface{} `json: "paras"`
}
