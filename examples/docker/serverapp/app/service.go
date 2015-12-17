// Echo server implementation taken from http://loige.co/simple-echo-server-written-in-go-dockerized/
//
// Set ENV var export DYNAMODB_LOCAL_ENDPOINT=http://192.168.99.101:8000
package main

import (
    "fmt"
    "net"
    "os"
    "strconv"
    "bytes"
    "net/url"
    "github.com/trustedhousesitters/roster"
)

const (
    CONN_PORT = "3333"
    CONN_TYPE = "tcp"
)

// Handles incoming requests.
func handleRequest(conn net.Conn) {
    // Make a buffer to hold incoming data.
    buf := make([]byte, 1024)

    // Read the incoming connection into the buffer.
    reqLen, err := conn.Read(buf)
    if err != nil {
        fmt.Println("Error reading:", err.Error())
    }

    // Builds the message.
    message := "Hi, I received your message! It was "
    message += strconv.Itoa(reqLen)
    message += " bytes long and that's what it said: \""
    n := bytes.Index(buf, []byte{0})
    message += string(buf[:n-1])
    message += "\" ! Honestly I have no clue about what to do with your messages, so Bye Bye!\n"

    // Write the message in the connection channel.
    conn.Write([]byte(message));
    // Close the connection when you're done with it.
    conn.Close()
}

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
        fmt.Println("Error getting local IP address:", err.Error())
        os.Exit(1)
    }

    endpoint := &url.URL{Scheme: CONN_TYPE, Host: CONN_HOST + CONN_PORT}

    service,err := client.Register("echo",endpoint.String())
    if err != nil {
        fmt.Println("Error registering:", err.Error())
        os.Exit(1)
    }

    defer service.Unregister()

    // Listen for incoming connections.
    l, err := net.Listen(CONN_TYPE, ":"+CONN_PORT)
    if err != nil {
        fmt.Println("Error listening:", err.Error())
        os.Exit(1)
    }
    // Close the listener when the application closes.
    defer l.Close()
    fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
    for {
        // Listen for an incoming connection.
        conn, err := l.Accept()
        if err != nil {
            fmt.Println("Error accepting: ", err.Error())
            os.Exit(1)
        }

        //logs an incoming message
        fmt.Printf("Received message %s -> %s \n", conn.RemoteAddr(), conn.LocalAddr())

        // Handle connections in a new goroutine.
        go handleRequest(conn)
    }
}
