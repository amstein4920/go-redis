package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type SavedData struct {
	value     string
	expirary  *int
	savedTime time.Time
}

func CommandsSwitch(commands []string, savedDataMap map[string]SavedData) string {
	switch strings.ToUpper(commands[0]) {
	case "PING":
		return "+PONG\r\n"
	case "ECHO":
		return fmt.Sprintf("$%v\r\n%v\r\n", len(commands[1]), commands[1])
	case "SET":
		return setCommand(commands, savedDataMap)
	case "GET":
		getCommand(commands, savedDataMap)
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
		expiraryValue, err := strconv.Atoi(commands[4])
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
	if savedDataForKey.value != "" && *savedDataForKey.expirary < time.Now().Compare(savedDataForKey.savedTime) {
		return fmt.Sprintf("$%v\r\n%v\r\n", len(savedDataForKey.value), savedDataForKey.value)
	}
	return "$-1\r\n"
}
