package main

import (
	"fmt"
	"net"
	"os"
	_"errors"
	"bufio"
	"strings"
	"time"
	_"golang.org/x/sync/semaphore"
	"strconv"
)

var SERVER_PORT = 9999
var server_name = ""

type Client struct{
	NickName string
	num string
	requests []string
	in_con net.Conn
}

type Objects struct{
	hashtable map[string]int
	locktable map[string]string
}

var transToclient map[string]string
var clientToTrans map[string]string
var data Data
var transToData = map[string]map[string]string{}
//var readLocks map[string]string
var locks map[string]*Lock

// make objects global

// handle request
func listener(cli *Client, inbox *Box){
	reader := bufio.NewReader(cli.in_con)
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
		m := Msg{}
		m.Parse(s)
		m.from = cli
		if (m.GetType() != "BEGIN") {
			m.transId = clientToTrans[m.from.num]
		}
		inbox.enqueue(m)
		// msg := strings.Split(s, " ")
	}
}


func startListening(inbox *Box, port int) {
	// fmt.Println("Started Listening on " + port)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		// handle error
		fmt.Printf("# [ERROR] %s", err)
	}
	c := 0
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			fmt.Printf("# [ERROR] %s", err)
		}
		cli := Client{in_con:conn}
		go listener(&cli, inbox)
		c++
	}
}


func set(m Msg) {
	// m.transId = clientToTrans[m.from.num]
	if _, ok := locks[m.obj]; !ok {
		// defining lock
		locks[m.obj] = &Lock{}
		locks[m.obj].init()
	}
	// if already a writer
	if locks[m.obj].isWriter(m.from.num) {
		// fmt.Println("updating write")
		transToData[m.transId][m.obj] = m.val
		m.from.SendMsg("OK")
		return
	}
	if ! locks[m.obj].lockWriter(m.from.num) {
		m.from.SendMsg("ABORTED")
		return
	}
	// fmt.Println("got write lock")
	transToData[m.transId][m.obj] = m.val
	// fmt.Println("transid: %s, obj:%s, val:%s",m.transId,m.obj,transToData[m.transId][m.obj])
	m.from.SendMsg("OK")
	return
}



func get(m Msg) {


	// if lock for the object is not defined
	if _, ok := locks[m.obj]; !ok {
		// defining lock
		locks[m.obj] = &Lock{}
		locks[m.obj].init()
	}
	// if has write lock
	if locks[m.obj].isWriter(m.from.num) {
		val, _ := transToData[m.transId][m.obj]
		s := fmt.Sprintf("%s.%s %s", server_name, m.obj, val)
		m.from.SendMsg(s)
		return
	}
	// fmt.Println("is writer not")
	// in case the object doesnt exists
	if _,ok := data.isPresent(m.obj); !ok {
		m.from.SendMsg("NOT FOUND")
		return
	}

	// if already a reader of this object
	if locks[m.obj].isReader(m.from.num) {
		val,_ := data.isPresent(m.obj)
		s := fmt.Sprintf("%s.%s %s", server_name, m.obj, val)
		m.from.SendMsg(s)
		return
	}

	// blocking for lock
	if ! locks[m.obj].lockReader(m.from.num){
		m.from.SendMsg("ABORTED")
		return
	}
	// lock required
	val,_ := data.isPresent(m.obj)
	s := fmt.Sprintf("%s.%s %s", server_name, m.obj, val)
	m.from.SendMsg(s)
	// wait to acquire the lock
	// fmt.Println("GET done")
	return
}

func main(){
	if len(os.Args) != 3 {
		fmt.Println("Need server id as argument. One of {A,B,C,D,E}      portnum")
	}
	server_name = os.Args[1]
	// fmt.Printf(server_name)
	inbox := Box{}

	//var connected_members map[string]*Node
	transToclient = make(map[string]string)
	clientToTrans = make(map[string]string)
	data = Data{}
	data.init()
	transToData = map[string]map[string]string{}
	// readLocks = make(map[string]string)
	locks = make(map[string]*Lock)
	data.update("x","100")
	// fmt.Println("Started Listening on\n")
	a,_ := strconv.Atoi(os.Args[2])
	go startListening(&inbox, a)
	for {
			m, err := inbox.pop()
			//fmt.Println("Got msg\n")
			// sleeping if no message in inbox
			if err != nil {
				time.Sleep(10)
				continue
			}
			switch m.GetType() {
			// client telling the server its ID
			case "ID":
				m.from.num = m.id
				break
			case "BEGIN":
				// if no transaction running from this client
				// fmt.Println("BEGIN Received")
				// fmt.Println(m.transId + " id")
				if _, ok := clientToTrans[m.from.num];!ok {
					clientToTrans[m.from.num] = m.transId
					transToclient[m.transId] = m.from.num

					transToData[m.transId] = map[string]string{}
					// transToData[m.transId]["x"] = "100"
					for _,v := range locks {
						v.kill.Off(m.from.num)
					}
					m.from.SendMsg("OK")
				} else {
					m.from.SendMsg("NOT_OK")
				}
				break
			case "GET":
				// fmt.Println("GET Received")
				go get(m)
				break
			case "SET":
				// fmt.Println("SET Received")
				go set(m)
				break
			case "PRECOMMIT":
				// do some kind of checking
				m.from.SendMsg("OK")
			case "COMMIT":
				// commiting all the data
				for v,k := range transToData[m.transId] {
					fmt.Printf("Commiting %s:%s\n",v,k)
					data.update(v,k)
				}
				for _,v := range locks {
					v.removeReader(m.from.num)
					v.removeWriter(m.from.num)
				}
				delete(transToData, m.transId)
				delete(transToclient, m.transId)
				delete(clientToTrans, m.from.num)
				// fmt.Println("Sending commit ok")
				m.from.SendMsg("COMMIT OK")
				break
			case "ABORT":
				for _,v := range locks {
					v.kill.On(m.from.num)
					v.removeReader(m.from.num)
					v.removeWriter(m.from.num)
				}
				delete(transToData, m.transId)
				delete(transToclient, m.transId)
				delete(clientToTrans, m.from.num)
				m.from.SendMsg("ABORTED")

				break
				// abort everything and clear all
			default:
				break
			}
	}
	return
}
