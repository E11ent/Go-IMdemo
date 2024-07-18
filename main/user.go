package main

import (
	"fmt"
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	addr := conn.RemoteAddr().String()
	name := conn.RemoteAddr().String()
	user := &User{
		Name: name,
		Addr: addr,
		C:    make(chan string),
		conn: conn,
	}
	go user.Listen()

	return user
}

func (user *User) Listen() {
	for {
		msg := <-user.C

		user.conn.Write([]byte(msg + "\r\n"))
		fmt.Println("msg:=", msg)
	}
}
