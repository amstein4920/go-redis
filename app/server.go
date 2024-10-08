package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

var portFlag *int
var replicaFlag *string

func init() {
	portFlag = flag.Int("port", 6379, "port for connection")
	replicaFlag = flag.String("replicaof", "", "leader server to replicate")
}

func main() {
	flag.Parse()

	if *replicaFlag != "" {
		processReplConn()
	}

	address := fmt.Sprintf("0.0.0.0:%v", *portFlag)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		fmt.Printf("Failed to bind to port %d", *portFlag)
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
	if *replicaFlag != "" {
		roleString = "slave"
	}

	//Calculation comes later, currently there is no ability to actually link to
	replicationOffset := strconv.FormatInt(0, 10)

	return ReplicationData{
		Role:               roleString,
		Master_replid:      "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
		Master_repl_offset: replicationOffset,
	}
}

func processReplConn() {
	leaderAddress := strings.Split(*replicaFlag, " ")

	port, err := strconv.Atoi(leaderAddress[1])
	if err != nil {
		fmt.Println("Invalid leader/master port provided")
		os.Exit(1)
	}

	address := fmt.Sprintf("%s:%d", leaderAddress[0], port)
	leaderConn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Failed to connected to leader/master server")
		os.Exit(1)
	}

	_, err = sendOverConn(leaderConn, "*1\r\n$4\r\nPING\r\n")
	if err != nil {
		fmt.Println("No Ping response received")
		return
	}

	_, err = sendOverConn(leaderConn,
		fmt.Sprintf("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n%d\r\n", *portFlag))
	if err != nil {
		fmt.Println("No REPLCONF response for port")
		return
	}
	_, err = sendOverConn(leaderConn, "*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n")
	if err != nil {
		fmt.Println("No REPLCONF response for psync")
		return
	}

	_, err = sendOverConn(leaderConn, "*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n")
	if err != nil {
		fmt.Println("No Psync response")
	}
}

func sendOverConn(leaderConn net.Conn, request string) (string, error) {
	response := make([]byte, 1024)
	fmt.Fprint(leaderConn, request)
	reader := bufio.NewReader(leaderConn)
	reader.Read(response)
	if string(response) == "" {
		return "", errors.New("no response from leader/master")
	}
	return string(response), nil
}
