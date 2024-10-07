package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
)

var portFlag *int
var replicaFlag *bool

func init() {
	portFlag = flag.Int("port", 6379, "port value")
	replicaFlag = flag.Bool("replicaof", false, "is replica flag")
}

func main() {
	flag.Parse()

	address := fmt.Sprintf("0.0.0.0:%v", *portFlag)
	listener, err := net.Listen("tcp", address)
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
	savedDataMap := make(map[string]SavedData)
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
				commands = append(commands, token)
			}
			if arrSize == 0 {
				break
			}
		}
		if len(commands) == 0 {
			break
		}

		replicationData := setReplicationData()

		response := CommandsSwitch(commands, savedDataMap, replicationData)

		_, err := conn.Write([]byte(response))
		if err != nil {
			fmt.Printf("[#%d] Error writing response: %v\n", id, err.Error())
			break
		}
	}

	fmt.Printf("#%d Client closing\n", id)
}

func setReplicationData() ReplicationData {
	roleString := "master"
	if *replicaFlag {
		roleString = "slave"
	}

	replicationOffset := strconv.FormatInt(0, 10)

	return ReplicationData{
		Role:               roleString,
		Master_replid:      "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		Master_repl_offset: replicationOffset,
	}
}
