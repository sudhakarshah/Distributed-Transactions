package main

import (
	"fmt"
	"net"
	"os"
	_"errors"
	"bufio"
	"strings"
)

// need map to server
func cli(server_map map[string]net.Conn)error{
	reader := bufio.NewReader(os.Stdin)
	// read input continuously
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		// broadcast the message to all servers
		handle_command(server_map, text)
	}
	return nil
}

// send PRE-COMMIT

func handle_command(server_map map[string]net.Conn, comm string){
	s := strings.Split(comm, " ")
	switch s[0]{
	case "BEGIN":
		broadcast(server_map, comm)
		resp := listen_all_response(server_map)
		if all_responses_are(resp, "OK"){
			fmt.Println("OK")
		}
	case "SET":
		if unicast_wrapper(server_map, s[1], comm) < 1{
			break
		}
		if listen_single_response(server_map, s[1]) != "OK"{
			fmt.Println("# COULD NOT SET. THIS SHOULD NOT HAPPEN")
		} else {
			fmt.Println("OK")
		}
	case "GET":
		if unicast_wrapper(server_map, s[1], comm) < 1{
			break
		}
	case "COMMIT":
		broadcast(server_map, "PRECOMMIT")
		resp := listen_all_response(server_map)
		if response_has(resp, "ABORTED"){
			broadcast(server_map, "ABORT")
			if !all_responses_are(resp, "ABORTED"){
				fmt.Println("# ABORT DID NOT WENT THROUGH AFTER BAD COMMIT. THIS SHOULD NOT HAPPEN!")
			} else {
				fmt.Println("ABORTED")
			}// all servers should respond with ABORTED
		}else{
			broadcast(server_map, "COMMIT")
			if !all_responses_are(resp, "COMMIT OK"){
				fmt.Println("# CHECKED FOR COMMIT BUT DID NOT ACTUALLY COMMIT. THIS SHOULD NOT HAPPEN!")
			} else {
				fmt.Println("COMMIT OK")
			}// all servers should respond with ABORTED
		}
	case "ABORT":
		broadcast(server_map, comm)
		resp := listen_all_response(server_map)
		if all_responses_are(resp, "ABORTED"){
			fmt.Println("ABORTED")
		} // all servers should respond with ABORTED
	default:
		fmt.Printf("# UNKOWN COMMAND: %s\n", comm)
	}

}

func broadcast(server_map map[string]net.Conn, comm string){
	for k, v := range server_map{
		unicast(v, comm, k)
	}
}

func unicast(server net.Conn, comm string, server_name string){
	_, err := fmt.Fprintf(server, fmt.Sprintf("%s\n", comm))
	if err != nil {
		fmt.Printf("# SENDING TO %s ERROR: %s\n", server_name, err)
	}
}

func unicast_wrapper(server_map map[string]net.Conn, serv_obj string, comm string)int{
	server_name := strings.Split(serv_obj, ".")[0]
	val, ok := server_map[server_name]
	if ok {
		unicast(val, comm, server_name)
		return 1
	}
	fmt.Printf("# SERVER NAMED \"%s\" DOES NOT EXIST\n", server_name)
	return -1
}

func listen_all_response(server_map map[string]net.Conn)map[string]string{
	responses := make(map[string]string)
	for k, v := range server_map{
		s, err := bufio.NewReader(v).ReadString('\n')
		if err != nil {
			fmt.Printf("# Failed listening \n", k)
			fmt.Printf("# ERROR: %s\n", err)
			return nil
		}
		responses[k] = s
	}
	return responses
}

func listen_single_response(server_map map[string]net.Conn, serv_obj string)string{
	s, err := bufio.NewReader(server_map[serv_obj]).ReadString('\n')
	if err != nil {
		fmt.Printf("# Failed listening \n", serv_obj)
		fmt.Printf("# ERROR: %s\n", err)
		return "# ERROR"
	}
	return s
}


func response_has(responses map[string]string, query string) bool{
	for _,v := range responses{
		if v == query{
			return true
		}
	}
	return false
}

func all_responses_are(responses map[string]string, query string) bool{
	for _,v := range responses{
		if v != query{
			return false
		}
	}
	return true
}
