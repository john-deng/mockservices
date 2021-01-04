package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"hidevops.io/hiboot/pkg/log"
)

func main() {
	//arguments := os.Args
	//if len(arguments) == 1 {
	//	fmt.Println("Please provide host:port.")
	//	return
	//}

	CONNECT := "localhost:8585" //arguments[1]
	conn, err := net.Dial("tcp", CONNECT)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		num, err := fmt.Fprintf(conn, text+"\n")
		if err == nil {
			log.Debugf("%v read", num)
			message, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Print("->: " + message)
			if strings.TrimSpace(string(text)) == "STOP" {
				fmt.Println("TCP client exiting...")
				return
			}
		}

	}
}
