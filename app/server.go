package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	fmt.Println("Listening on: ", listener.Addr().String())

	for id := 1; ; id++ {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go serveClient(id, conn)
	}
}

func serveClient(id int, conn net.Conn) {
	defer conn.Close()
	for {
		scanner := bufio.NewScanner(conn)

		commands := []string{}
		var arrSize, strSize int
		for scanner.Scan() {
			token := scanner.Text()
			switch token[0] {
			case '*':
				arrSize, _ = strconv.Atoi(token[1:])
			case '$':
				strSize, _ = strconv.Atoi(token[1:])
			default:
				if len(token) != strSize {
					fmt.Println("Improper string size")
					break
				}
				arrSize--
				strSize = 0
				fmt.Println(token)
				commands = append(commands, token)
			}
			if arrSize == 0 {
				break
			}
		}
		if len(commands) == 0 {
			break
		}

		for _, command := range commands {
			fmt.Println(command)
		}

		var response string
		switch strings.ToUpper(commands[0]) {
		case "PING":
			response = "+PONG\r\n"
		case "ECHO":
			response = fmt.Sprintf("$%v\r\n%v\r\n", len(commands[1]), commands[1])
		}

		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Printf("[#%d] Error writing response: %v\n", id, err.Error())
			break
		}
	}

	fmt.Printf("[#%d] Client closing\n", id)
}
