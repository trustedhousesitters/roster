// Echo server - based on https://github.com/hirolovesbeer/golang-memo/blob/master/echo-server.go
// Set ENV var export DYNAMODB_LOCAL_ENDPOINT=http://192.168.99.101:8000

package main

import (
    "net"
    "log"
    "bufio"
    "os"
    "net/url"
    "os/signal"
    "syscall"
    "github.com/trustedhousesitters/roster"
)

const (
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
)

func main() {

    //Check to see if running Dynamodb Locally or in AWS
    var client *roster.Client
    if dle := os.Getenv("DYNAMODB_LOCAL_ENDPOINT"); dle != "" {
        client = roster.NewClient(roster.LocalConfig{Endpoint: dle})
    } else {
        client = roster.NewClient(roster.WebServiceConfig{})
    }

    CONN_HOST,err := client.GetLocalIP()
    if err != nil {
        log.Fatalln(err)
    }

    endpoint := &url.URL{Scheme: CONN_TYPE, Host: CONN_HOST + CONN_PORT}

    service,err := client.Register("echo",endpoint.String())
    if err != nil {
        log.Fatalln(err)
    }

    ln, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
    if err != nil {
        log.Fatalln(err)
        return
    }
    defer ln.Close()
    log.Printf("Listening on %s:%s\n", CONN_HOST,CONN_PORT)

    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigs
        log.Println("Unregistering service")
        service.Unregister()
        os.Exit(1)
    }()

    for {
        conn, _ := ln.Accept()
        go func(conn net.Conn) {
            //logs an incoming message
            log.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())

            reader := bufio.NewReader(conn)
            for {
                line, _, e := reader.ReadLine()
                if e != nil {
                    conn.Close()
                    return
                }
                log.Printf("> %s\n", line)

            }
        }(conn)
    }
}
