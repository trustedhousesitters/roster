// Echo server - based on https://github.com/hirolovesbeer/golang-memo/blob/master/echo-client.go
// Set ENV var export DYNAMODB_LOCAL_ENDPOINT=http://192.168.99.101:8000

package main

import (
    "net"
    "log"
    "os"
    "net/url"
    "time"
    "github.com/trustedhousesitters/roster"
)

func main() {
    client := roster.NewClient(roster.ClientConfig{})

    // Send a message every 3 seconds
    for {
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

        time.Sleep(3  * time.Second)
    }

}
