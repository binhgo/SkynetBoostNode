package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func InitServer(TcpPort string, insertQueue chan *Node, updateQueue chan *Node, getQueue chan net.Conn) {

	fmt.Printf("Server is listening on port %s\n", TcpPort)

	tcpAddr, err := net.ResolveTCPAddr("tcp4", TcpPort)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		conn, err := listener.Accept()

		log.Printf("Accepted conn from %s\n", conn.RemoteAddr().String())

		if err != nil {
			continue
		}

		ip, port := getIPAndPort(conn.RemoteAddr())

		// en queue new node, another go routine will save it to database later
		node := NewNode(ip, port, "6789", time.Now(), time.Now(), true)
		insertQueue <- node

		go handleClient(conn, updateQueue, getQueue)
	}
}

func handleClient(conn net.Conn, updateQueue chan *Node, getQueue chan net.Conn) {
	defer conn.Close()

	data := make([]byte, 1024)

	for {
		len, err := conn.Read(data)

		if err != nil {
			fmt.Println(err)
			return
		}

		if len != 0 {
			fmt.Printf("Receive data from client: %s\n", string(data[:len]))

			command := string(data[:len])
			// check data here and send back ACK to client
			// 1. update active status
			// 2. get all connected nodes
			if command == "update" {
				ip, port := getIPAndPort(conn.RemoteAddr())
				node := NewNode(ip, port, "6789", time.Now(), time.Now(), true)

				// another go routine will check this channel and push the update into database
				updateQueue <- node

			} else if command == "get" {
				// GET queue
				// another go routine will check this channel and push data to clients
				getQueue <- conn

			} else {
				continue
			}
		}
	}
}

func SendDataToClient(conn net.Conn, data string) {
	if conn != nil {
		_, err := conn.Write([]byte(data))
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func getIPAndPort(addr net.Addr) (IP string, Port string) {
	s := strings.Split(addr.String(), ":")
	return s[0], s[1]
}
