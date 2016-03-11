package main

import (
    "os"
    "net"
    "fmt"
    "log"
    "time"
    "html"
    "bufio"
    "strconv"
    "syscall"
    "net/http"
    "math/rand"
    "io/ioutil"
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
    port, _ := strconv.Atoi(envPort)

    if port == 0 {
        port = 8080 + index
    }

    urlDependency := os.Getenv(fmt.Sprint("SAVA_RUNNER_HTTP_DEPENDENCY_URL", index))

    fmt.Println("Listening on HTTP port: ", port)

    if len(urlDependency) > 0 {
        fmt.Println("Will proxy requests to: ", urlDependency)
    }

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        path := html.EscapeString(r.URL.Path)[1:]

        if len(urlDependency) == 0 {
            httpResponse(Response{id, runtime, port, path}, w)
        } else {
            resp, err := http.Get(fmt.Sprint(urlDependency, "/", path))
            if err != nil {
                fmt.Println("Error connection to dependency ", urlDependency, " - ", err.Error())
                httpResponse(Response{id, runtime, port, ""}, w)
            } else {
                defer resp.Body.Close()
                body, err := ioutil.ReadAll(resp.Body)
                if err != nil {
                    fmt.Println("Error getting content from dependency ", urlDependency, " - ", err.Error())
                    httpResponse(Response{id, runtime, port, ""}, w)
                } else {
                    w.Header().Set("Content-Type", "application/json")
                    w.Write(body)
                }
            }
        }
    })

    log.Fatal(http.ListenAndServe(":" + strconv.Itoa(port), nil))
}

func httpResponse(response Response, w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

func serveTcp(id, runtime string, index int) {

    fmt.Println("Starting TCP server: ", index)

    envPort := os.Getenv(fmt.Sprint("SAVA_RUNNER_TCP_PORT", index))
    port, _ := strconv.Atoi(envPort)

    if port == 0 {
        port = 8090 + index
    }

    urlDependency := os.Getenv(fmt.Sprint("SAVA_RUNNER_TCP_DEPENDENCY_URL", index))

    if len(urlDependency) > 0 {
        fmt.Println("Will proxy requests to: ", urlDependency)
    }

    fmt.Println("Listening on TCP port: ", port)
    listener, err := net.Listen("tcp", fmt.Sprint(":", port))
    if err != nil {
        log.Fatal("Error listening: ", err.Error())
    }

    defer listener.Close()

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
        } else {
            go handleTcpRequest(id, runtime, port, conn, urlDependency)
        }
    }
}

func handleTcpRequest(id, runtime string, port int, conn net.Conn, urlDependency string) {
    buffer := make([]byte, 1024)
    length, err := conn.Read(buffer)
    if err != nil {
        fmt.Println("Error reading: ", err.Error())
    } else {
        if len(urlDependency) > 0 {
            client, err := net.Dial("tcp", urlDependency)
            if err != nil {
                fmt.Println("Error reading: ", urlDependency, " - ", err.Error())
                conn.Write([]byte(fmt.Sprint("{\"id\":\"", id, "\",\"runtime\":\"", runtime, "\",\"port\":", port, ",\"request\":\"\"}")))
            } else {
                client.Write(buffer[:length])
                read := make([]byte, 1024)
                len, err := bufio.NewReader(client).Read(read)
                if err != nil {
                    fmt.Println("Error reading: ", urlDependency, " - ", err.Error())
                    conn.Write([]byte(fmt.Sprint("{\"id\":\"", id, "\",\"runtime\":\"", runtime, "\",\"port\":", port, ",\"request\":\"\"}")))
                } else {
                    conn.Write(read[:len])
                }
            }
        } else {
            request := string(buffer[:length])
            conn.Write([]byte(fmt.Sprint("{\"id\":\"", id, "\",\"runtime\":\"", runtime, "\",\"port\":", port, ",\"request\":\"", request, "\"}")))
        }
    }
    conn.Close()
}
