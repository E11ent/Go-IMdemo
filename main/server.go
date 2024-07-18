package main

import (
	"fmt"
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
	fmt.Println("链接创建成功")
	server.mapLock.Lock()
	user := NewUser(conn)
	server.onlineMap[user.Name] = user
	server.mapLock.Unlock()

	server.boardCast(user, "is online")
	fmt.Println("消息发送成功")

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
