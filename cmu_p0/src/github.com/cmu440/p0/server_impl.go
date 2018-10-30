// Implementation of a MultiEchoServer. Students should write their code in this file.

package p0

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

const (
	defaultHost = "localhost"
	defaultPort = 9999
)

type echoClient struct {
	conn net.Conn
	ch   chan string
}

type multiEchoServer struct {
	// TODO: implement this!
	host        string
	port        int
	clientsChan chan map[int]*echoClient
	stop        chan bool
	ln          net.Listener
}

// New creates and returns (but does not start) a new MultiEchoServer.
func New() MultiEchoServer {
	mes := &multiEchoServer{
		host:        defaultHost,
		port:        defaultPort,
		clientsChan: make(chan map[int]*echoClient, 1),
		stop:        make(chan bool),
	}
	return MultiEchoServer(mes)
}

func (mes *multiEchoServer) Start(port int) error {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	mes.ln = ln
	if err != nil {
		fmt.Println("Start error: ", err)
		return err
	}

	mes.clientsChan <- map[int]*echoClient{}

	go mes.serve()

	return nil
}

func (mes *multiEchoServer) Close() {
	// TODO: implement this!

	clients := <-mes.clientsChan
	for _, client := range clients {
		client.conn.Close()
	}
	mes.ln.Close()

	close(mes.stop)
}

func (mes *multiEchoServer) Count() int {
	// TODO: implement this!
	clients := <-mes.clientsChan
	count := len(clients)
	mes.clientsChan <- clients
	return count
}

func (mes *multiEchoServer) serve() {

	msgChan := make(chan string)
	go handleMessage(msgChan, mes.clientsChan)

	i := 0
	for {
		conn, err := mes.ln.Accept()
		if err != nil {
			select {
			case <-mes.stop:
				fmt.Println("Exit serve")
				return
			default:
				// mes.stop <- false
				// fmt.Println("Accept error: ", err)
			}
			fmt.Println("Accept error: ", err)
			continue
		}
		go handleConnection(i, conn, mes.clientsChan, msgChan)
		i++
	}
}

// TODO: add additional methods/functions below!
func handleConnection(i int, conn net.Conn, clientsChan chan map[int]*echoClient, msgChan chan string) {
	fmt.Printf("Client %d: %v <-> %v\n", i, conn.LocalAddr(), conn.RemoteAddr())

	// Close connection when this function ends
	defer func() {
		fmt.Printf("%d: closed\n", i)
		clients := <-clientsChan
		delete(clients, i)
		clientsChan <- clients
		conn.Close()
	}()

	clients := <-clientsChan
	client := &echoClient{conn, make(chan string, 100)}
	clients[i] = client
	clientsChan <- clients

	go sendBack(client)

	bufReader := bufio.NewReader(conn)

	for {
		bytes, err := bufReader.ReadString('\n')
		if err != nil {
			break
		}
		msgChan <- bytes
	}
}

// keep sending message
func handleMessage(msgChan <-chan string, clientsChan chan map[int]*echoClient) {
	for {
		msg := <-msgChan
		clients := <-clientsChan
		for _, client := range clients {
			// sendBack(client.conn, msg)
			select {
			case client.ch <- msg:
			default:
			}
		}
		clientsChan <- clients
	}
}

func sendBack(client *echoClient) {
	for {
		msg := <-client.ch
		_, err := client.conn.Write([]byte(msg))
		if err != nil {
			return
		}
	}
}
