package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	listen()
}
func handleData(data string) {
	db, _ := sql.Open("mysql", "root:@/test")
	statement, err := db.Prepare("INSERT INTO test.buttonLog(ButtonData,DateStamp) values(?,?)")
	if err != nil {
		log.Fatal(err)
	}
	current_time := time.Now()
	res, err := statement.Exec(data, current_time)
	fmt.Println(res)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
func handleConnection(connection net.Conn) {
	buf := make([]byte, 1024)
	reqLen, err := connection.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
		connection.Write([]byte("NOK\n"))
		return
	}
	fmt.Println(reqLen)
	connBuf := bufio.NewReader(connection)
	request, err := connBuf.ReadString('\n')
	fmt.Println(err)
	fmt.Println(request)

	connection.Write([]byte("OK\n"))
	go handleData(request)
	connection.Close()
}
func listen() {
	ln, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("Error binding to socket")
		os.Exit(1)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error receiving data")
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}
