package main

import (
	"net"
	"fmt"
)

var NUM_TO_SERVER_NAME = map[int]string{6:"A",7:"B",8:"C",9:"D",10:"E"}
var SERVER_PORT = 9999

// connects to all servers
func get_servers(vm_num string)map[string]net.Conn{
	server_map := make(map[string]net.Conn)
	for i := 5; i<10 ;i++{
		server_name := NUM_TO_SERVER_NAME[i+1]
		url := fmt.Sprintf("sp19-cs425-g%s-%02d.cs.illinois.edu", vm_num, i+1)
		//fmt.Printf("%s == %s\n", url, NUM_TO_SERVER_NAME[i+1])
		conn := dial_server(url)
		server_map[server_name] = conn
	}
	return server_map
}

func dial_server(url string)net.Conn{
	//fmt.Printf("%s == %s\n", url, NUM_TO_SERVER_NAME[i+1])
	full_url := fmt.Sprintf("%s:%d", url, SERVER_PORT)
	conn, err := net.Dial("tcp", full_url)
	if err != nil {
		return nil
	}
	return conn
}

func main(){
	//server_map := get_servers("62")
	conn := dial_server("127.0.0.1")
	server_map := map[string]net.Conn{"A": conn}
	cli(server_map)
}
