package main

import (
	"fmt"
	"net"
	"os"
	_"errors"
	"bufio"
	"strings"
	"math/rand"
	"time"
	"sync"
)

var curr_trans = ""
var counterLock sync.Mutex
var commandsHandling = 0
// need map to server
func cli(server_map map[string]net.Conn)error {
	reader := bufio.NewReader(os.Stdin)
	// read input continuously from user
	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)
		// broadcast the message to all servers
		counterLock.Lock()
		// if a command is already running and the new command is not abort then continue
		if (commandsHandling != 0 && strings.Split(text, " ")[0] != "ABORT") {
			counterLock.Unlock()
			continue
		}
		commandsHandling += 1
		counterLock.Unlock()
		go handle_command(server_map, text)
	}
	return nil
}

// send PRE-COMMIT

func handle_command(server_map map[string]net.Conn, comm string) {
	s := strings.Split(comm, " ")
	switch s[0] {
		case "BEGIN":
			// adding unique transaction ID to the string

			curr_trans = RandomString(10)
			comm += " " + RandomString(10)
			broadcast(server_map, comm)
			resp := listen_all_response(server_map)
			if all_responses_are(resp, "OK") {
				fmt.Println("OK")
			}
			break
		case "SET":

			if unicast_wrapper(server_map, s[1], comm) < 1 {
				break
			}
			if listen_single_response(server_map, strings.Split(s[1], ".")[0]) != "OK" {
				fmt.Println("# COULD NOT SET. THIS SHOULD NOT HAPPEN")
			} else {
				fmt.Println("OK")
			}
		case "GET":
			if unicast_wrapper(server_map, s[1], comm) < 1 {
				break
			}
			// fmt.Println("going to listen\n")
			s := listen_single_response(server_map, strings.Split(s[1], ".")[0])
			if (s == "NOT FOUND") {
				// aborting the transaction
				handle_command(server_map, "ABORT")
				break
			}
			ss := strings.Split(s, " ")
			fmt.Printf("%s = %s\n",ss[0],ss[1])
			break
			// print out the value here
		case "COMMIT":
			broadcast(server_map, "PRECOMMIT")
			resp := listen_all_response(server_map)
			if response_has(resp, "ABORTED") {
				broadcast(server_map, "ABORT")
				// probably need to call this also
				resp = listen_all_response(server_map)
				//
				if !all_responses_are(resp, "ABORTED") {
					fmt.Println("# ABORT DID NOT WENT THROUGH AFTER BAD COMMIT. THIS SHOULD NOT HAPPEN!")
				} else {
					fmt.Println("ABORTED")
				}// all servers should respond with ABORTED
			} else {
				broadcast(server_map, "COMMIT")
				// probably need to call this also
				resp := listen_all_response(server_map)
				//
				if !all_responses_are(resp, "COMMIT OK") {
					fmt.Println("# CHECKED FOR COMMIT BUT DID NOT ACTUALLY COMMIT. THIS SHOULD NOT HAPPEN!")
				} else {
					fmt.Println("COMMIT OK")
				}// all servers should respond with COMMIT OK
			}
		case "ABORT":
			broadcast(server_map, comm)
			resp := listen_all_response(server_map)
			if all_responses_are(resp, "ABORTED"){
				fmt.Println("ABORTED")
			} // all servers should respond with ABORTED
			curr_trans = ""
		default:
			fmt.Printf("# UNKOWN COMMAND: %s\n", comm)
			break
		}
		counterLock.Lock()
		commandsHandling -= 1
		counterLock.Unlock()
		return

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

func listen_all_response(server_map map[string]net.Conn)map[string]string {
	responses := make(map[string]string)
	// listenting to responses of all the servers
	for k, v := range server_map {
		s, err := bufio.NewReader(v).ReadString('\n')
		if err != nil {
			fmt.Printf("# Failed listening \n", k)
			fmt.Printf("# ERROR: %s\n", err)
			return nil
		}
		fmt.Printf("# server %s replied : %s\n", k, s)
		responses[k] = strings.TrimSpace(s)
	}
	return responses
}

func listen_single_response(server_map map[string]net.Conn, serv_obj string)string {
	fmt.Printf("server_obj %s\n", serv_obj)
	s, err := bufio.NewReader(server_map[serv_obj]).ReadString('\n')
	if err != nil {
		fmt.Printf("# Failed listening \n", serv_obj)
		fmt.Printf("# ERROR: %s\n", err)
		return "# ERROR"
	}
	fmt.Printf("# server %s alone replied : %s\n", serv_obj, s)
	return strings.TrimSpace(s)
}


func response_has(responses map[string]string, query string) bool{
	for _,v := range responses {
		if v == query {
			return true
		}
	}
	return false
}

func all_responses_are(responses map[string]string, query string) bool{
	for k,v := range responses {
		if v != query {
			fmt.Printf("for server %s, %s != %s\n",k,v,query)
			return false
		}
	}
	return true
}

func RandomString(len int) string {
	rand.Seed(time.Now().UnixNano())
	bytes := make([]byte, len)
		for i := 0; i < len; i++ {
			bytes[i] = byte(65 + rand.Intn(25))  //A=65 and Z = 65+25
		}
	return string(bytes)
}
