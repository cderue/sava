package main

import (
    "os"
    "fmt"
    "flag"
    "strconv"
    "syscall"
    "os/signal"
    "github.com/gin-gonic/gin"
)

func main() {

    // Waiter keeps the program from exiting instantly.
    waiter := make(chan bool)

    cleanup := func() {
        fmt.Println("Exit.")
    }

    // Catch a CTR+C exits so the cleanup routine is called.
    c := make(chan os.Signal, 1)
    signal.Notify(c, os.Interrupt)
    signal.Notify(c, syscall.SIGTERM)

    go func() {
        <-c
        cleanup()
        os.Exit(1)
    }()

    defer cleanup()

    count := count()

    for i := 1; i <= count; i++ {
        go serve(i)
    }

    waiter <- true
}

func count() int {
    envCount := os.Getenv("SAVA_RUNNER_COUNT")
    count, _ := strconv.Atoi(envCount)

    if count == 0 {
        count = 1
    }

    return count
}

func serve(index int) {

    fmt.Println("Starting server: ", index)

    envPort := os.Getenv(fmt.Sprint("SAVA_RUNNER_PORT", index))
    envPortNumber, _ := strconv.Atoi(envPort)

    if envPortNumber == 0 {
        envPortNumber = 8080 + index
    }

    arg := fmt.Sprint("port", index)
    description := fmt.Sprint("Sets the port ", index, " to listen on")

    port := flag.Int(arg, envPortNumber, description)

    flag.Parse()

    r := gin.Default()
    r.GET("/", func(c *gin.Context) {
        c.String(200, fmt.Sprint(index))
    })

    fmt.Println("Listening on port: ", *port)
    r.Run(":" + strconv.Itoa(*port))
}
