package main

import (
        "log"
        "net/http"
        "os"
	"encoding/json"

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
        log.Printf("Watching for completed analysis reports for delivery");
        socket, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
                log.Print("upgrade:", err)
                return
        }
        defer socket.Close()
        for {
                messages := getMessages(healthAnalysisQueue)
                for _, thisMessage := range messages {
			var jsonPandaHealthAnalysis string = *thisMessage.Body
                        socket.WriteMessage(websocket.TextMessage, []byte(jsonPandaHealthAnalysis))
                        if err != nil {
                                log.Println("Failed to write message to websocket:", err)
                                break
                        } else {
				deleteMessage(thisMessage, healthAnalysisQueue)
				var pandaHealthAnalysis PandaHealthAnalysis
				json.Unmarshal([]byte(jsonPandaHealthAnalysis), &pandaHealthAnalysis)
				if err == nil {
					log.Printf("Successfully delivered health analysis results for panda %s", pandaHealthAnalysis.Name)
				} else {
					log.Printf("Delivered health results for some panda, but could not determine name %s", jsonPandaHealthAnalysis)

				}
			}
                }
        }
}
