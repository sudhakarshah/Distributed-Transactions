package main
import(
	"sync"
	"errors"
	_"net"
	_"encoding/json"
	"fmt"
	_"time"
	_"bufio"
	"strings"
)

type Data struct {
	data map[string]string
	mux sync.Mutex
}

func (d *Data)init() {
	d.data = make(map[string]string)
}

func (d *Data)update(k string, v string) {
	d.mux.Lock()
	d.data[k] = v
	d.mux.Unlock()
}
func (d *Data)isPresent(k string)(string, bool) {
	d.mux.Lock()
	defer d.mux.Unlock()
	v,ok := d.data[k]
	return v,ok

}

type Msg struct{
	Type string
	obj string
	val string
	transId string
	from *Client
	id string
}

type Box struct{
	messages []Msg
	mux sync.Mutex
}

func (m *Msg)Parse(s string)int {
	tokens := strings.Split(strings.TrimSpace(s), " ")
	m.Type = tokens[0]
	switch m.GetType() {
	case "BEGIN":
		m.transId = tokens[1]
		return 1
	case "GET":
		m.obj = strings.Split(tokens[1], ".")[1]
		return 1
	case "SET":
		m.obj = strings.Split(tokens[1], ".")[1]
		m.val = tokens[2]
		return 1
	case "ID":
		m.id = tokens[1]
	default:
		return -1
	}
	return -1
}


func (m * Msg)GetType()string{
	return m.Type
}


func (in*Box) enqueue(m Msg){
	in.mux.Lock()
	in.messages = append(in.messages, m)
	in.mux.Unlock()
}

func (in*Box) pop()(Msg,error){
	var output Msg
	var err error
	//fmt.Printf("pop called\n")
	in.mux.Lock()
	if len(in.messages) != 0{
		output = in.messages[0]
		in.messages = append(in.messages[:0], in.messages[1:]...)
	} else {
		err = errors.New("The inbox is empty")
		in.mux.Unlock()
		return output, err
	}
	in.mux.Unlock()
	return output, err
}

func (cli *Client)SendMsg(s string)int {
	// fmt.Println("sending " + s+"\n")
	_, err := fmt.Fprintf(cli.in_con, s+"\n")
	if err != nil {
		fmt.Printf("# Failed sending to node %s\n", cli.num)
		fmt.Printf("# ERROR: %s\n", err)
		return 0
	}
	return 1
}
