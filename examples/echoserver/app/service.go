// Echo server - based on https://github.com/hirolovesbeer/golang-memo/blob/master/echo-server.go

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
    client := roster.NewClient(roster.ClientConfig{})

    CONN_HOST,err := client.GetLocalIP()
    if err != nil {
        log.Fatalln(err)
    }

    endpoint := &url.URL{Scheme: CONN_TYPE, Host: CONN_HOST + ":" + CONN_PORT}

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
                line, err := reader.ReadBytes('\n')
                if err != nil { // EOF, or worse
                    break
                }
                conn.Write(line)
            }
        }(conn)
    }
}
