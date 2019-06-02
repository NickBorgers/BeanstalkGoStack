package main

import (
	"log"
        "encoding/json"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	for {
		messages := getMessages(requestQueue)
		for _, thisMessage := range messages {
			processPandaDataRequest(thisMessage)
		}
	}
}

func processPandaDataRequest(message *sqs.Message) {
	pandaHealthData := getHealthDataForPandaByName(message.Body)
	jsonPandaHealthData,err := json.Marshal(pandaHealthData)
	if err == nil {
		sendMessage(jsonPandaHealthData, healthDataQueue)
		log.Printf("Retrieved and sent along data for requested panda %s", message.Body)
		deleteMessage(message, requestQueue)
	} else {
        	log.Printf("Could not build data for requested panda %s", message.Body)
	}
}
