package main

import (
    "log"
    "net/http"
    "os"
)

func main() {
	for {
		messages := getMessages(requestQueue)
		for _, thisMessage := range messages {
			processPandaDataRequest(thisMessage)
		}
	}
}

func processPandaDataRequest(message *Message) {
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
