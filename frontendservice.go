package main

import (
    "log"
    "net/http"
    "os"
)

func main() {
	port := os.Getenv("PORT")
        if port == "" {
	        port = "5000"
        }

        http.HandleFunc("/pandas/", pandaRequestHandler);

        log.Printf("Listening on port %s\n\n", port)
        http.ListenAndServe(":"+port, nil)
}

func pandaRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusAccepted)
               	log.Printf("Received GET: %s\n", r.URL.Path)
                sendMessage(r.URL.Path, "testqueue")
        } else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
    		log.Printf("Rejected %s: %s\n", r.Method, r.URL.Path)
        }
}

