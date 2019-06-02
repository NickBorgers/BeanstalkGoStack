package main

import (
        "log"
        "encoding/json"
        "github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
        for {
                messages := getMessages(healthDataQueue)
                for _, thisMessage := range messages {
                        processPandaAnalysisRequest(thisMessage)
                }
        }
}

func processPandaAnalysisRequest(message *sqs.Message) {
        var jsonPandaHealthData string = *message.Body
        var pandaHealthData PandaHealthData

        err := json.Unmarshal([]byte(jsonPandaHealthData), &pandaHealthData)
        if err == nil {
                pandaHealthAnalysis := analyzePandaHealthData(pandaHealthData)

                jsonPandaHealthAnalysis,err := json.Marshal(pandaHealthAnalysis)
                if err == nil {
                        sendMessage(string(jsonPandaHealthAnalysis), healthAnalysisQueue)
                        log.Printf("Analyzed and sent along analysis for health of panda %s", pandaHealthData.Name)
                        deleteMessage(message, healthDataQueue)
                } else {
                        log.Printf("Could not analyze health data for panda %s", pandaHealthData.Name)
                }
        } else {
                log.Printf("Could not parse health data for panda: %s", jsonPandaHealthData)
        }
}

func analyzePandaHealthData(healthData PandaHealthData) PandaHealthAnalysis {
        
        var maximumRemainingLifeExpectancy int

        if healthData.Status != "dead" {

                maximumRemainingLifeExpectancy = int(pandaMaxLife - healthData.Age)
        
                lifeExpectancyImpacts := make([]int, len(healthData.HealthIndicators))

                for index, thisHealthIndicator := range healthData.HealthIndicators {
                        var expectancyImpact = thisHealthIndicator.LifeExpectancyImpact*float32(maximumRemainingLifeExpectancy)
                        lifeExpectancyImpacts[index] = int(expectancyImpact)
                }

                for _, thisLifeExpectancyImpact := range lifeExpectancyImpacts {
                        maximumRemainingLifeExpectancy += thisLifeExpectancyImpact
                }

                if maximumRemainingLifeExpectancy <= 0 {
                        maximumRemainingLifeExpectancy = 0
                }

        } else {
                maximumRemainingLifeExpectancy = int(0)
        }

        var healthAnalysis PandaHealthAnalysis
        // Copy over data fields from received
        healthAnalysis.Name = healthData.Name
        healthAnalysis.Status = healthData.Status
        healthAnalysis.Age = healthData.Age
        healthAnalysis.HealthIndicators = healthData.HealthIndicators
        
        healthAnalysis.ExpectedAgeAtDeath = uint(int(healthAnalysis.Age) + maximumRemainingLifeExpectancy)
        healthAnalysis.ExpectedYearsOfMortalityRemaining = uint(healthAnalysis.ExpectedAgeAtDeath) - healthAnalysis.Age

        return healthAnalysis
}
