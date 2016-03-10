package main

import (
    "os"
    "fmt"
    "log"
    "time"
    "flag"
    "html"
    "strconv"
    "syscall"
    "net/http"
    "math/rand"
    "os/signal"
    "encoding/json"
)

var hex = []rune("0123456789ABCDEF")

type Response struct {
    Id      string `json:"id"`
    Runtime string `json:"runtime"`
    Port    int    `json:"port"`
    Path    string `json:"path"`
}

func main() {

    rand.Seed(time.Now().UnixNano())

    var id = os.Getenv("SAVA_RUNNER_ID")
    var runtime = random(16)

    fmt.Println("Id     : ", id)
    fmt.Println("Runtime: ", runtime)

    // Waiter keeps the program from exiting instantly.
    waiter := make(chan bool)

    cleanup := func() {}

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
        go serve(id, runtime, i)
    }

    waiter <- true
}

func random(n int) string {
    b := make([]rune, n)
    for i := range b {
        b[i] = hex[rand.Intn(len(hex))]
    }
    return string(b)
}

func count() int {
    envCount := os.Getenv("SAVA_RUNNER_COUNT")
    count, _ := strconv.Atoi(envCount)

    if count == 0 {
        count = 1
    }

    return count
}

func serve(id, runtime string, index int) {

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

    fmt.Println("Listening on port: ", *port)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        path := html.EscapeString(r.URL.Path)[1:]
        response := Response{id, runtime, *port, path}

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })

    log.Fatal(http.ListenAndServe(":" + strconv.Itoa(*port), nil))
}
