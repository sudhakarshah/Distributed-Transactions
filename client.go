package main

import (
	"net"
	"fmt"
	"os"
)

var NUM_TO_SERVER_NAME = map[int]string{6:"A",7:"B",8:"C",9:"D",10:"E"}
var SERVER_PORT = 9999
var my_id = ""
// var port = 5000
// connects to all servers
func get_servers(vm_num string)map[string]net.Conn{
	server_map := make(map[string]net.Conn)
	fmt.Println(vm_num)
	for i := 5; i<10 ;i++{
		server_name := NUM_TO_SERVER_NAME[i+1]
		url := fmt.Sprintf("sp19-cs425-g%s-%02d.cs.illinois.edu", vm_num, i+1)
		//url := "127.0.0.1"
		//fmt.Printf("%s == %s\n", url, NUM_TO_SERVER_NAME[i+1])
		conn := dial_server(url,SERVER_PORT)
		server_map[server_name] = conn
	}
	return server_map
}

func dial_server(url string,port int)net.Conn{
	//fmt.Printf("%s == %s\n", url, NUM_TO_SERVER_NAME[i+1])
	full_url := fmt.Sprintf("%s:%d", url, port)
	conn, err := net.Dial("tcp", full_url)
	if err != nil {
		return nil
	}
	_, err = fmt.Fprintf(conn, fmt.Sprintf("ID %s\n", my_id))
	if err != nil {
		return nil
	}
	return conn
}

func main(){

	if len(os.Args) != 2 {
		fmt.Println("Need client id as argument: 0-10")
	}
	my_id = os.Args[1]
	fmt.Printf("sending\n")
	server_map := get_servers("62")
	//conn := dial_server("127.0.0.1",5000)
	// conn2 := dial_server("127.0.0.1",5001)
	fmt.Printf("connected\n")
	// server_map := map[string]net.Conn{"A": conn}
	// server_map["B"] = conn2
	cli(server_map)
}
