package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
)

const (
	defaultHost = "localhost"
	defaultPort = 9999
)

// To test your server implementation, you might find it helpful to implement a
// simple 'client runner' program. The program could be very simple, as long as
// it is able to connect with and send messages to your server and is able to
// read and print out the server's echoed response to standard output. Whether or
// not you add any code to this file will not affect your grade.
func main() {
	conn, err := net.Dial("tcp", defaultHost+":"+strconv.Itoa(defaultPort))
	if err != nil {
		fmt.Println("ERROR Dial: ", err)
		os.Exit(-1)
	}

	ch := make(chan string)

	go sendMessage(conn, ch)

	go receiveMessage(conn)

	for {
		var input string
		fmt.Scanln(&input)

		ch <- input
	}
}

func sendMessage(conn net.Conn, ch <-chan string) {
	bufWriter := bufio.NewWriter(conn)
	for {
		select {
		case msg := <-ch:
			_, err := bufWriter.WriteString(msg + "\n")
			bufWriter.Flush()
			if err != nil {
				fmt.Println("ERROR WRITE: ", err)
				os.Exit(-1)
			}
		}
	}
}

func receiveMessage(conn net.Conn) {
	bufReader := bufio.NewReader(conn)
	for {
		line, err := bufReader.ReadString('\n')
		if err != nil {
			break
		}
		fmt.Print(line)
	}
}
