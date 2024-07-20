package main

import (
	"net"
	"strings"
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
	if msg == "ls" {
		newMsg := "Now is online:\r\n"
		user.server.mapLock.RLock()
		for _, u := range user.server.onlineMap {
			newMsg += "[" + u.Addr + "]" + u.Name
		}
		user.server.mapLock.RUnlock()
		user.C <- newMsg
	} else if len(msg) >= 7 && msg[:7] == "rename " {
		newName := msg[7:]
		_, ok := user.server.onlineMap[newName]
		if ok {
			user.C <- "fail : new name exists\r\n"
		} else {
			user.server.mapLock.Lock()
			_, ok2 := user.server.onlineMap[newName]
			if ok2 {
				user.C <- "fail : new name exists\r\n"
				return
			}
			delete(user.server.onlineMap, user.Name)
			user.server.onlineMap[newName] = user
			user.Name = newName
			defer user.server.mapLock.Unlock()
		}
		user.C <- "rename to " + user.Name + " successfully"
	} else if len(msg) >= 3 && msg[:3] == "to " {
		ordArray := strings.Split(msg, " ")
		targetName := ordArray[1]
		var msgContext string
		for index, str := range ordArray {
			if index > 1 {
				msgContext = msgContext + str + " "
			}
		}
		target := user.server.onlineMap[targetName]
		target.C <- user.Name + ": " + msgContext
		user.C <- "to " + target.Name + ": " + msgContext
	} else {
		user.server.boardCast(user, msg)
	}
}

func (user *User) Listen() {
	for {
		msg := <-user.C
		user.conn.Write([]byte(msg + "\r\n"))
	}
}
