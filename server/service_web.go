package main

import(
	"log"
	"fmt"
	"net/http"
	"html/template"
	"os"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"strconv"
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
            "server" : "server"+sep,
            "static" : "client"+sep+"static"+sep,
            "view"   : "client"+sep+"view"+sep,
        },
    }
    htmls = []string{
				config.path["view"]+"base.html",
				config.path["view"]+"index.html",
		}
)

type Data struct{
		Date string `json:"date"`
		Item string `json:"item"`
		Payer string `json:"payer"`
		State string `json:"state"`
		Reimburse string `json:"reimburse"`
		Income string `json:"income"`
		Outcome string `json:"outcome"`
}
func addHandler(w http.ResponseWriter, r *http.Request) {
		tmpl,err := template.ParseFiles(htmls[0],htmls[1])
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
				//fmt.Println("recv add!")
				addInfo(data)
				c.Emit("added","")
		})
		//-----------------------------------------------------
		http.Handle("/socket.io/", server)

}

func startWeb(){
	InitSocket()
	http.HandleFunc("/", addHandler) // static version
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(config.path["static"]))))
	fmt.Println("Server Started on Port ", config.port)
	var portstr = os.Getenv("PORT")
	if portstr!=""{
		log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
	}else{
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.port), nil))
	}
	
	
}
