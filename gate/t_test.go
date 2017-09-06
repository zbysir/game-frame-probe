package main

import (
	"testing"
	"github.com/gorilla/websocket"
	"github.com/bysir-zl/bygo/log"
	"strconv"
	"time"
)

func BenchmarkConn(b *testing.B) {
	server := "47.94.204.137"
	d := websocket.DefaultDialer
	for i := 0; i < b.N; i++ {
		conn, _, err := d.Dial("ws://"+server+":8081/", nil)
		if err != nil {
			log.ErrorT("test", err)
			continue
		}

		iS := strconv.Itoa(i)
		conn.WriteMessage(websocket.BinaryMessage, []byte(`{"cmd":1,"body":"`+iS+`"}`))
		conn.WriteMessage(websocket.BinaryMessage, []byte(`{"cmd":100,"body":"`+iS+`"}`))
	}

}

func TestOnline(t *testing.T) {
	server := "47.94.204.137"
	d := websocket.DefaultDialer
	conns := []*websocket.Conn{}
	for i := 0; i < 300; i++ {
		conn, _, err := d.Dial("ws://"+server+":8081/", nil)
		if err != nil {
			log.ErrorT("test", err)
			continue
		}
		conns = append(conns, conn)
		iS := strconv.Itoa(i)
		conn.WriteMessage(websocket.BinaryMessage, []byte(`{"cmd":1,"body":"`+iS+`"}`))
		conn.WriteMessage(websocket.BinaryMessage, []byte(`{"cmd":100,"body":"`+iS+`"}`))
	}
	log.InfoT("test", "over")

	for {
		for _, c := range conns {
			c.WriteMessage(websocket.BinaryMessage, []byte(`{"cmd":110,"body":"sb"}`))
		}
		time.Sleep(1 * time.Second)
	}
}

func TestConn2(t *testing.T) {
	server := "47.94.204.137"
	d := websocket.DefaultDialer

	conn, _, err := d.Dial("ws://"+server+":8081/", nil)
	if err != nil {
		log.ErrorT("test", err)
		return
	}

	conn.WriteMessage(websocket.BinaryMessage, []byte(`{"cmd":1,"body":"2"}`))
	conn.WriteMessage(websocket.BinaryMessage, []byte(`{"cmd":100,"body":"2"}`))

	go func() {
		for range time.Tick(2 * time.Second) {
			conn.WriteMessage(websocket.BinaryMessage, []byte(`{"cmd":110,"body":"2"}`))
		}
	}()

	for {
		_, body, err := conn.ReadMessage()
		if err != nil {
			log.ErrorT("test", err)
			break
		}
		log.Info(string(body))
	}
}

func TestConn1(t *testing.T) {
	server := "47.94.204.137"
	//server := "localhost"
	d := websocket.DefaultDialer

	conn, _, err := d.Dial("ws://"+server+":8081/", nil)
	if err != nil {
		log.ErrorT("test", err)
		return
	}

	conn.WriteMessage(websocket.BinaryMessage, []byte(`{"cmd":1,"body":"1"}`))
	conn.WriteMessage(websocket.BinaryMessage, []byte(`{"cmd":100,"body":"1"}`))

	go func() {
		for range time.Tick(2 * time.Second) {
			conn.WriteMessage(websocket.BinaryMessage, []byte(`{"cmd":110,"body":"1"}`))
		}
	}()

	for {
		_, body, err := conn.ReadMessage()
		if err != nil {
			log.ErrorT("test", err)
			break
		}
		log.Info(string(body))
	}
}
