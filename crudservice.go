package main

import (
    "io/ioutil"
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
	if r.Method == "POST" {
        	if buf, err := ioutil.ReadAll(r.Body); err == nil {
                	log.Printf("Received POST: %s\n", string(buf))
                }
        } else {
                if buf, err := ioutil.ReadAll(r.Body); err == nil {
    			log.Printf("Received GET: %s\n", string(buf))
                }
        }
}
