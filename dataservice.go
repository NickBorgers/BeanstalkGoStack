package main

import (
        "log"
        "encoding/json"
        "github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
	// For as long as we run
        for {
		// Read names of pandas to pull data for from the queue
                messages := getMessages(requestQueue)
                for _, thisMessage := range messages {
                        processPandaDataRequest(thisMessage)
                }
        }
}

func processPandaDataRequest(message *sqs.Message) {
        var requestedPandaName string = *message.Body
	// Lookup this panda's health information in our sophisticated PHI system
        pandaHealthData := getHealthDataForPandaByName(requestedPandaName)
	// JSON encode it
        jsonPandaHealthData,err := json.Marshal(pandaHealthData)
        if err == nil {
		// And shove it on down the next queue to the analysis service (or not, who knows)
                sendMessage(string(jsonPandaHealthData), healthDataQueue)
                log.Printf("Retrieved and sent along data for requested panda %s", requestedPandaName)
		// Delete this data request so we only proccess it one time
                deleteMessage(message, requestQueue)
        } else {
                log.Printf("Could not build data for requested panda %s", requestedPandaName)
        }
}
