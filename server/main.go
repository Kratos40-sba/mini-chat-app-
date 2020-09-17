package main

import (
	"bufio"
	"log"
	"net"
)

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	openConnections = make(map[net.Conn]bool)
	newConnection   = make(chan net.Conn)
	deadConnection  = make(chan net.Conn)
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	handleError(err)
	defer l.Close()
	go func() {
		for {
			conn, err := l.Accept()
			handleError(err)
			openConnections[conn] = true
			newConnection <- conn
		}
	}()
	for {
		select {
		case conn := <-newConnection:
			go broadcastConnection(conn)
		case conn := <-deadConnection:
			for c := range openConnections {
				if c == conn {
					break
				}
			}
			delete(openConnections, conn)

		}
	}

}
func broadcastConnection(conn net.Conn) {
	for {
		reader := bufio.NewReader(conn)
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		for c := range openConnections {
			if c != conn {
				c.Write([]byte(msg))
			}
		}

	}
	deadConnection <- conn
}
