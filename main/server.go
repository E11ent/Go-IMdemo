package main

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
)

type Server struct {
	ip   string
	port int

	onlineMap map[string]*User
	mapLock   sync.RWMutex

	message chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{ip: ip, port: port, onlineMap: make(map[string]*User), message: make(chan string)}
}

func (server *Server) listenMessage() {
	for {
		msg := <-server.message
		server.mapLock.RLock()
		for _, user := range server.onlineMap {
			user.C <- msg
		}
		server.mapLock.RUnlock()
	}
}

func (server *Server) boardCast(user *User, msg string) {
	msg = "[" + user.Addr + "]" + user.Name + ":" + msg
	server.message <- msg
}

func (server *Server) handleConnection(conn net.Conn) {

	user := NewUser(conn, server)
	user.Online()

	go func() {
		var msg string
		for {
			buf := make([]byte, 2)
			n, err := conn.Read(buf)
			fmt.Println("n==", n)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				return
			}
			str := string(buf)
			if str != "\r\n" {
				msg += str
			} else {
				user.DoMessage(msg)
				msg = ""
			}

		}
	}()

}

func (server *Server) Run() {
	listener, err := net.Listen("tcp", server.ip+":"+strconv.Itoa(server.port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	go server.listenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go server.handleConnection(conn)
	}
}
