package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Second))
	for {
		readBuffer := make([]byte, 1024)
		_, err := conn.Read(readBuffer)
		if err != nil {
			fmt.Println("Error reading data from connection: ", err.Error())
			os.Exit(1)
		}
		_, err = conn.Write([]byte("+PONG\r\n"))
		if err != nil {
			fmt.Println("Error writing PONG back  to client: ", err.Error())
		}
		conn.SetDeadline(time.Now().Add(time.Second))
	}
}
