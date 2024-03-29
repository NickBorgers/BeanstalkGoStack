package main

import (
        "log"
        "net/http"
        "os"
        "encoding/json"

        "github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool)
var upgrader = websocket.Upgrader{}

func main() {
        port := os.Getenv("PORT")
        if port == "" {
                port = "5000"
        }

        // Register redirect handler for getting folks to the HTML page
        http.HandleFunc("/", redirectToHomePage)

        // Register handlers for incoming HTTP requests
        http.HandleFunc("/pandas/", pandaRequestHandler);
        http.HandleFunc("/pandas/healthReports", handleWebsocketConnections);

        // Startup goroutine which watches for analysis result messages
        go retrieveAndSendAnalysisReports()
        log.Printf("Watching for completed analysis reports for delivery");

        // Starting listening for HTTP requests
        log.Printf("Listening on port %s\n\n", port)
        http.ListenAndServe(":"+port, nil)
}

func redirectToHomePage(w http.ResponseWriter, r *http.Request) {
        http.Redirect(w, r, "/index.html", 302)
}

func pandaRequestHandler(w http.ResponseWriter, r *http.Request) {
        // Serve requests by kicking of a synchronous request for panda health information
        if r.Method == "GET" {
                w.WriteHeader(http.StatusAccepted)
                log.Printf("Received GET: %s\n", r.URL.Path)
                sendMessage(r.URL.Path, requestQueue)
        } else {
        // Reject anything but a GET because we cannot serve them
                http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
                    log.Printf("Rejected %s: %s\n", r.Method, r.URL.Path)
        }
}

func handleWebsocketConnections(w http.ResponseWriter, r *http.Request) {
        // Upgrade initial GET request to a websocket
        socket, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
                log.Fatal(err)
        }

        // Register our new client
        clients[socket] = true

        log.Printf("New websocket client registered");
}

func retrieveAndSendAnalysisReports() {
        // For as long as we run
        for {
                // Check for new messages on the analysis queue
                messages := getMessages(healthAnalysisQueue)
                // For each of them we get back
                for _, thisMessage := range messages {
                        // Get the string content
                        var jsonPandaHealthAnalysis string = *thisMessage.Body

                        var atLeastOneSuccessfulDelivery = false

                        // And attempt delivery of it to all connected clients
                        for client := range clients {
                                err := client.WriteMessage(websocket.TextMessage, []byte(jsonPandaHealthAnalysis))
                                // If a client is unavailable
                                if err != nil {
                                        log.Printf("Failed to send message to a client: %v", err)
                                        // Close the connection on our side and delete reference to it
                                        client.Close()
                                        delete(clients, client)
                                } else {
                                        // We got through to at least one client, hooray!
                                        atLeastOneSuccessfulDelivery = true
                                }
                        }

                        // Attempt to determine what panda this message was about for logging success or failure
                        var pandaHealthAnalysis PandaHealthAnalysis
                        err := json.Unmarshal([]byte(jsonPandaHealthAnalysis), &pandaHealthAnalysis)

                        // If any client received our message
                        if atLeastOneSuccessfulDelivery {
                                // Delete delivered message
                                deleteMessage(thisMessage, healthAnalysisQueue)
        
                                if err == nil {
                                        log.Printf("Successfully delivered health analysis results for panda %s", pandaHealthAnalysis.Name)
                                } else {
                                        log.Printf("Delivered health results for some panda, but could not determine name %s", jsonPandaHealthAnalysis)
                                }
                        } else {
                                if err == nil {
                                        log.Printf("Failed to deliver to any client the health analysis results for panda %s", pandaHealthAnalysis.Name)
                                } else {
                                        log.Printf("Failed to deliver to any client the health analysis results for report that coudld not be parsed %s", jsonPandaHealthAnalysis)
                                }
                        }
                }
        }

}
