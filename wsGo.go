package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "185.251.90.101:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options
var Connections []*websocket.Conn

func echo(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	Connections = append(Connections, c)
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		for _, Connection:=range Connections{
			if(Connection != c) {
				err = Connection.WriteMessage(mt, message)
				if err != nil {
					log.Println("write:", err)
					break
				}
			}
		}

	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}