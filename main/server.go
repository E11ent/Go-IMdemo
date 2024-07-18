package main

import (
	"fmt"
	"net"
	"strconv"
)

type Server struct {
	ip   string
	port int
}

func NewServer(ip string, port int) *Server {
	return &Server{ip: ip, port: port}
}

func (server *Server) handleConnection(conn net.Conn) {
	fmt.Println("链接创建成功")
}

func (server *Server) Run() {
	listener, err := net.Listen("tcp", server.ip+":"+strconv.Itoa(server.port))
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go server.handleConnection(conn)
	}
}
