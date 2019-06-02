package main

import (
        "log"
        "net/http"
        "os"

	"github.com/gorilla/websocket"
)

func main() {
	port := os.Getenv("PORT")
        if port == "" {
	        port = "5000"
        }

        http.HandleFunc("/pandas/", pandaRequestHandler);
	http.HandleFunc("/pandas/healthReports", pandaHealthAnalysisReportsHandler);

        log.Printf("Listening on port %s\n\n", port)
        http.ListenAndServe(":"+port, nil)
}

func pandaRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusAccepted)
               	log.Printf("Received GET: %s\n", r.URL.Path)
                sendMessage(r.URL.Path, requestQueue)
        } else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
    		log.Printf("Rejected %s: %s\n", r.Method, r.URL.Path)
        }
}

var upgrader = websocket.Upgrader{}

func pandaHealthAnalysisReportsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
