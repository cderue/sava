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
    "net"
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

    httpCount := countHttp()

    for i := 1; i <= httpCount; i++ {
        go serveHttp(id, runtime, i)
    }

    tcpCount := countTcp()

    for i := 1; i <= tcpCount; i++ {
        go serveTcp(id, runtime, i)
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

func countHttp() int {
    return count("SAVA_RUNNER_HTTP_COUNT")
}

func countTcp() int {
    return count("SAVA_RUNNER_TCP_COUNT")
}

func count(env string) int {
    envCount := os.Getenv(env)
    count, _ := strconv.Atoi(envCount)
    if count == 0 {
        count = 1
    }
    return count
}

func serveHttp(id, runtime string, index int) {

    fmt.Println("Starting HTTP server: ", index)

    envPort := os.Getenv(fmt.Sprint("SAVA_RUNNER_HTTP_PORT", index))
    envPortNumber, _ := strconv.Atoi(envPort)

    if envPortNumber == 0 {
        envPortNumber = 8080 + index
    }

    arg := fmt.Sprint("http-port-", index)
    description := fmt.Sprint("Sets the HTTP port ", index, " to listen on")

    port := flag.Int(arg, envPortNumber, description)

    flag.Parse()

    fmt.Println("Listening on HTTP port: ", *port)

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        path := html.EscapeString(r.URL.Path)[1:]
        response := Response{id, runtime, *port, path}

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    })

    log.Fatal(http.ListenAndServe(":" + strconv.Itoa(*port), nil))
}

func serveTcp(id, runtime string, index int) {

    fmt.Println("Starting TCP server: ", index)

    envPort := os.Getenv(fmt.Sprint("SAVA_RUNNER_TCP_PORT", index))
    envPortNumber, _ := strconv.Atoi(envPort)

    if envPortNumber == 0 {
        envPortNumber = 8090 + index
    }

    arg := fmt.Sprint("tcp-port-", index)
    description := fmt.Sprint("Sets the TCP port ", index, " to listen on")

    port := flag.Int(arg, envPortNumber, description)

    flag.Parse()

    fmt.Println("Listening on TCP port: ", *port)

    listener, err := net.Listen("tcp", fmt.Sprint("localhost:", *port))
    if err != nil {
        log.Fatal("Error listening: ", err.Error())
    }

    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
        } else {
            go handleTcpRequest(id, runtime, *port, conn)
        }
    }
}

func handleTcpRequest(id, runtime string, port int, conn net.Conn) {
    buffer := make([]byte, 1024)
    length, err := conn.Read(buffer)
    if err != nil {
        fmt.Println("Error reading:", err.Error())
    } else {
        request := string(buffer[:length])
        conn.Write([]byte(fmt.Sprint("{\"id\":\"", id, "\",\"runtime\":\"", runtime, "\",\"port\":", port, ",\"request\":\"", request, "\"}")))
    }
    conn.Close()
}
