// Echo server - based on https://github.com/hirolovesbeer/golang-memo/blob/master/echo-client.go
// Set ENV var export DYNAMODB_LOCAL_ENDPOINT=http://192.168.99.101:8000

package main

import (
    "net"
    "log"
    "os"
    "net/url"
    "github.com/trustedhousesitters/roster"
)

func main() {

    //Check to see if running Dynamodb Locally or in AWS
    var client *roster.Client
    if dle := os.Getenv("DYNAMODB_LOCAL_ENDPOINT"); dle != "" {
        client = roster.NewClient(roster.LocalConfig{Endpoint: dle})
    } else {
        client = roster.NewClient(roster.WebServiceConfig{})
    }

    service,err := client.Discover("echo")
    if err != nil {
        log.Fatalln(err)
    }

    u, err := url.Parse(service.Endpoint)
	if err != nil {
		log.Fatal(err)
	}

    strEcho := "Hello echo service"

    tcpAddr, err := net.ResolveTCPAddr("tcp", u.Host)
    if err != nil {
        log.Fatal(err)
    }

    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        log.Fatal(err)
    }

    _, err = conn.Write([]byte(strEcho))
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Write to server = ", strEcho)

    reply := make([]byte, 1024)

    _, err = conn.Read(reply)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Reply from server=", string(reply))

    conn.Close()
}
