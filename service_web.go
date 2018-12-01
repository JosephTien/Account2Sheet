package main

import (
	"log"
	"fmt"
	"net/http"
	"html/template"
	"os"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"strconv"
	"strings"
)
type Config struct{
    port int
    path map[string]string
}
var(
    sep = string(os.PathSeparator)
    config = Config{
        port : 8080,
        path : map[string]string{
            "static" : "static"+sep,
            "templates" : "templates"+sep,
        },
    }
    templates = []string{
		config.path["templates"]+"base.tmpl.html",
		config.path["templates"]+"index.tmpl.html",
	}
)

type Data struct{
		SpreadsheetId  string `json:"spreadsheetId"`
		TableName string `json:"tableName"`
		Date string `json:"date"`
		Item string `json:"item"`
		State string `json:"state"`
		Payer string `json:"payer"`
		Receipt string `json:"receipt"`
		Reimburse string `json:"reimburse"`
		Income string `json:"income"`
		Outcome string `json:"outcome"`
}
func addHandler(w http.ResponseWriter, r *http.Request) {
		tmpl,err := template.ParseFiles(templates[0],templates[1])
				if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = tmpl.ExecuteTemplate(w, "base", nil)
		if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
		}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "helloworld")
}

func InitSocket(){
		server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
		//---------- OnConnection ---------------------------------------
		server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
			log.Println("Connected")
			c.Join("default")
		})

		//---------- OnDisconnection ------------------------------------
		server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
			log.Println("Disconnected")
		})
		//---------- OnAdd ------------------------------------
		server.On("add", func(c *gosocketio.Channel, data Data) {
				error := addInfo(data)
				if !error{
					c.Emit("added","")
				}else{
					c.Emit("error","")
				}
		})
		//---------- OnList ------------------------------------
		server.On("requirelist", func(c *gosocketio.Channel, data Data) {
			readList(data.SpreadsheetId, data.TableName)
			var sheetlist []string = make([]string, len(spreadsheets))
			i := 0
			for key, value := range spreadsheets {
				sheetlist[i]=fmt.Sprintf("%s/%s", key, value.gid)
				i++
			}
			c.Emit("list",strings.Join(sheetlist,"\\"))
		})
		//-----------------------------------------------------
		http.Handle("/socket.io/", server)
}

func startWeb(){
	InitSocket()
	http.HandleFunc("/", addHandler) // static version
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.path["static"]))))
	var portstr = os.Getenv("PORT")
	if portstr!=""{
		fmt.Println("Server Started on Port ", portstr)
		log.Fatal(http.ListenAndServe(":"+portstr, nil))
	}else{
		fmt.Println("Server Started on Port ", config.port)
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.port), nil))
	}
}
