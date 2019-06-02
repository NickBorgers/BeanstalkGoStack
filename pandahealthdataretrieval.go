package main

import (
        "hash/fnv"
)

func getHealthDataForPandaByName(name string) PandaHealthData {
        var pandaHealthData PandaHealthData
        pandaHealthData.Name = name

        var nameKey = hash(name)

        // Lookup if panda is alive
        var isAliveKey = nameKey % 7
        var status = "living"
        if isAliveKey < 2 {
                status = "dead"
        }
        pandaHealthData.Status = status

        // Lookup current age of panda
        var currentAge = nameKey % uint32(pandaMaxLife)
        pandaHealthData.Age = uint(currentAge)

        // Lookup health indicators
        var numberOfHealthIndicators = nameKey % 3 + 1
        healthIndicators := make([]HealthIndicator, numberOfHealthIndicators)
        for i := uint32(0); i< numberOfHealthIndicators; i++ {
                var thisHealthIndicator = getHealthIndicator(nameKey + nameKey*i)
                healthIndicators[i] = thisHealthIndicator
        }

        pandaHealthData.HealthIndicators = healthIndicators

        return pandaHealthData
}

func hash(s string) uint32 {
        h := fnv.New32a()
        h.Write([]byte(s))
        return h.Sum32()
}

func getHealthIndicator(key uint32) (HealthIndicator) {
        var healthIndicatorIndex = key % uint32(len(pandaHealthIndicators))
        return pandaHealthIndicators[healthIndicatorIndex]
}
