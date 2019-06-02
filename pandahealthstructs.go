package main

type PandaHealthData struct {
    Name string `json:"name"`
    Status string `json:"status"`
    Age uint `json:"age"`
    HealthIndicators []HealthIndicator `json:"healthIndicators"`
}

type PandaHealthAnalysis struct {
    Name string `json:"name"`
    Status string `json:"status"`
    Age uint `json:"age"`
    ExpectedAgeAtDeath uint `json:"expectedAgeAtDeath"`
    ExpectedYearsOfMortalityRemaining uint `json:"expectedYearsOfMortalityRemaining"`
    HealthIndicators []HealthIndicator `json:"healthIndicators"`
}

type HealthIndicator struct {
    Name string `json:"name"`
    LifeExpectancyImpact float32 `json:"lifeExpectancyImpact"`
}
