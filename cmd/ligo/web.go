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
		tp, mess, err := conn.ReadMessage()
		fmt.Println(tp)
		if err != nil {
			log.Println(err)
			return
		}
		part := strings.TrimSpace(ligo.StripComments(string(mess)))
		fmt.Printf("Message recieved : %v\n", []byte(part))
		message := ""
		v, err := vm.Eval(part)
		message = fmt.Sprint(v.Value)
		if v.Value == nil {
			message = "<i class=\"nilReturn\">nil</i>"
		}
		if err != nil {
			message = "<b class=\"error\">" + err.Error() + "</b>"
		}
		err = conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println(err)
			return
		}
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
