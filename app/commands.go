package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type SavedData struct {
	value     string
	expirary  *int64
	savedTime time.Time
}

// Format of field names nonstandard so as to easily match what the Redis protocol specifies
type ReplicationData struct {
	Role               string
	Master_replid      string
	Master_repl_offset string
}

func CommandsSwitch(commands []string, savedDataMap map[string]SavedData, replicationData ReplicationData) string {
	switch strings.ToUpper(commands[0]) {
	case "PING":
		return "+PONG\r\n"
	case "ECHO":
		return fmt.Sprintf("$%v\r\n%v\r\n", len(commands[1]), commands[1])
	case "SET":
		return setCommand(commands, savedDataMap)
	case "GET":
		return getCommand(commands, savedDataMap)
	case "INFO":
		return infoCommand(commands, replicationData)
	}
	return "$-1\r\n"
}

func setCommand(commands []string, savedDataMap map[string]SavedData) string {
	if len(commands) < 4 {
		savedDataMap[commands[1]] = SavedData{
			value: commands[2],
		}
		return "+OK\r\n"
	}

	if strings.ToUpper(commands[3]) == "PX" {
		expiraryValue, err := strconv.ParseInt(commands[4], 10, 64)
		if err != nil {
			fmt.Println("Improper Expirary Value")
			return "$-1\r\n"
		}
		savedDataMap[commands[1]] = SavedData{
			value:     commands[2],
			expirary:  &expiraryValue,
			savedTime: time.Now(),
		}
	}
	return "+OK\r\n"
}

func getCommand(commands []string, savedDataMap map[string]SavedData) string {
	savedDataForKey := savedDataMap[commands[1]]
	if savedDataForKey.value == "" {
		return "$-1\r\n"
	}
	if savedDataForKey.expirary == nil {
		return fmt.Sprintf("$%v\r\n%v\r\n", len(savedDataForKey.value), savedDataForKey.value)
	}
	if *savedDataForKey.expirary > time.Since(savedDataForKey.savedTime).Milliseconds() {
		return fmt.Sprintf("$%v\r\n%v\r\n", len(savedDataForKey.value), savedDataForKey.value)
	}
	return "$-1\r\n"
}

func infoCommand(commands []string, replicationData ReplicationData) string {
	switch commands[1] {
	case "replication":
		dataString := ""

		rdReflection := reflect.ValueOf(replicationData)
		typeOfRD := rdReflection.Type()
		for i := 0; i < rdReflection.NumField(); i++ {
			dataString = strings.Join([]string{dataString, strings.ToLower(typeOfRD.Field(i).Name) + ":" + rdReflection.Field(i).Interface().(string)}, " ")
		}
		return fmt.Sprintf("$%v\r\n%v\r\n", len(dataString), dataString)
	}
	return "$-1\r\n"
}
