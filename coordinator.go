package main

import (
	"fmt"
	"net"
	"os"
	_"errors"
	"bufio"
	_"strings"
)

var SERVER_PORT = 9999

type Client struct{
	NickName string
	requests []string
}

type Objects struct{
	hashtable map[string]int
	locktable map[string]string
}

// make objects global

// handle request
func listener(in_con net.Conn){
	reader := bufio.NewReader(in_con)
	// create local request buffer
	for {
		s, err := reader.ReadString('\n')
		fmt.Printf("# recieved string %s\n", s)
		if err != nil {
			fmt.Println("#Error in listening")
			fmt.Printf("# %s", err)
			return
		}
		s = strings.TrimSpace(s)
		msg := strings.Split(s, " ")
	}
}

func handle_command(server_map map[string]net.Conn, comm string) {

}

func startListening() {
	//fmt.Println("Started Listening on " + port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", SERVER_PORT))
	if err != nil {
		// handle error
		fmt.Printf("# [ERROR] %s", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Printf("# [ERROR] %s", err)
		}
		go listener(conn)
	}
}


func main(){
	if len(os.Args) != 2 {
		fmt.Println("Need server id as argument. One of {A,B,C,D,E}")
	}
	startListening()
	return
}
