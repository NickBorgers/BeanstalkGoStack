package main

import (
    "log"
    "net/http"
    "os"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sqs"
)

var svc *sqs.SQS

func main() {
	port := os.Getenv("PORT")
        if port == "" {
	        port = "5000"
        }

        sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc = sqs.New(sess)

        http.HandleFunc("/pandas/", pandaRequestHandler);

        log.Printf("Listening on port %s\n\n", port)
        http.ListenAndServe(":"+port, nil)
}

func pandaRequestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.WriteHeader(http.StatusAccepted)
               	log.Printf("Received GET: %s\n", r.URL.Path)
                sendMessage(r.URL.Path)
        } else {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
    		log.Printf("Rejected %s: %s\n", r.Method, r.URL.Path)
        }
}

func sendMessage(name string) (*sqs.SendMessageOutput) {
	// URL to our queue
        qURL := "https://sqs.us-east-1.amazonaws.com/181171223360/testqueue"

        result, err := svc.SendMessage(&sqs.SendMessageInput{
  		MessageBody: aws.String(name),
		QueueUrl:    &qURL,
	})

	if err != nil {
		log.Printf("Error", err)
		return nil
	}

	return result
}

