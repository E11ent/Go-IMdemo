package main

import (
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	addr := conn.RemoteAddr().String()
	name := conn.RemoteAddr().String()
	user := &User{
		Name:   name,
		Addr:   addr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.Listen()

	return user
}

func (user *User) Online() {
	user.server.mapLock.Lock()
	user.server.onlineMap[user.Name] = user
	user.server.mapLock.Unlock()
	user.DoMessage("is online")
}

func (user *User) Offline() {
	user.server.mapLock.Lock()
	delete(user.server.onlineMap, user.Name)
	user.server.mapLock.Unlock()
	user.DoMessage("is offline")
}

func (user *User) DoMessage(msg string) {
	user.server.boardCast(user, msg)
}

func (user *User) Listen() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg + "\r\n"))
	}
}
