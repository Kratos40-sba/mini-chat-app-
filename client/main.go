package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func main() {
	conn, err := net.Dial("tcp", "localhost:8080")
	handleError(err)
	defer conn.Close()
	fmt.Println("Enter your name : ")
	reader := bufio.NewReader(os.Stdin)
	username, err := reader.ReadString('\n')
	handleError(err)
	username = strings.Trim(username, "\r\n")
	welcom := fmt.Sprintf("welcome %s , to the chat \n", username)
	fmt.Println(welcom)
	// read cuncurntly
	// writing
	go read(conn)
	write(conn, username)

}
func write(connection net.Conn, username string) {
	for {
		reader := bufio.NewReader(os.Stdin)
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		msg = fmt.Sprintf("%s-> %s\n", username, strings.Trim(msg, "\r\n"))
		connection.Write([]byte(msg))
	}
}
func read(connection net.Conn) {
	for {
		reader := bufio.NewReader(connection)
		msg, err := reader.ReadString('\n')
		if err == io.EOF {
			connection.Close()
			fmt.Println("connection closed")
			os.Exit(0)
		}
		fmt.Println(msg)
		fmt.Println("******************************************")

	}
}
