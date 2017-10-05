package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aki237/ligo/pkg/ligo"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func serve(conn *websocket.Conn, vm *ligo.VM) {

	defer func() {
		if r := recover(); r != nil {
			message := "<b class=\"error\">" + fmt.Sprint(r) + "</b>"
			err := conn.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Println(err)
				conn.Close()
				return
			}
			go serve(conn, vm)
			return
		}
	}()

	for {
		_, mess, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		part := strings.TrimSpace(ligo.StripComments(string(mess)))
		message := ""
		v, err := vm.Eval(part)
		message = fmt.Sprint(v.Value)
		if err != nil {
			message = "<b class=\"error\">" + err.Error() + "</b>"
		}
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Message got : ", part, ", Message wrote : ", message)
	}
}

func handle(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	vm := ligo.NewVM()
	vm.Funcs["require"] = VMRequire
	vm.Funcs["load-plugin"] = VMDlLoad

	serve(conn, vm)
}

func runWeb() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(htmlPage))
	})

	http.HandleFunc("/ws", handle)
	http.ListenAndServe(":8080", nil)
}
