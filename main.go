package main

import (
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
)

type HomePageData struct {
	Title          string
	WelcomeMessage string
}

var wsUpgrader = websocket.Upgrader{}

func myMiddleware(handler http.Handler) http.Handler {
	log.Print("aa!")
	return handler
}

func main() {
	//homeTemplate := template.Must(template.ParseFiles("home.html"))

	//http.HandleFunc("/", homePageHandler(homeTemplate))
	//http.HandleFunc("/ws", handleWs())

	http.Handle("/", myMiddleware(http.FileServer(http.Dir("./next/out"))))

	http.HandleFunc("/api/", func(writer http.ResponseWriter, request *http.Request) {
		_, err := writer.Write([]byte("hi!"))
		if err != nil {
			log.Print(err)
		}
	})

	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		log.Panic(err)
	}
}

//type HTMLDir struct {
//	d http.Dir
//}
//
//func (d HTMLDir) Open(name string) (http.File, error) {
//	log.Print("OPEN CALLED LOL")
//	return d.d.Open(name + ".html")
//}

func handleWs() func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		log.Print("New client connected to the websocket.")
		conn, err := wsUpgrader.Upgrade(writer, request, nil)
		if err != nil {
			log.Print(err)
			return
		}
		defer conn.Close()

		for {
			messageType, data, _ := conn.ReadMessage()

			log.Printf("Received message: [%s]", data)

			err := conn.WriteMessage(messageType, []byte("Hi!"))
			if err != nil {
				log.Print(err)
				return
			}
		}

	}
}

func homePageHandler(homeTemplate *template.Template) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		homePageData := HomePageData{
			Title:          "Home page",
			WelcomeMessage: "Welcome from go template!",
		}

		err := homeTemplate.Execute(writer, homePageData)
		if err != nil {
			log.Panic(err)
		}
		log.Print("Served home page.")
	}
}
