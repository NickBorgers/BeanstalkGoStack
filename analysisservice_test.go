package main

import (
	"testing"
)

func TestAnalyzePandaHealthDataWithLivePanda(t *testing.T) {
	var testHealthData PandaHealthData
	testHealthData.Name = "TestPanda"
	testHealthData.Status = "living"
	testHealthData.Age = 12
	testHealthData.HealthIndicators = []HealthIndicator {
		HealthIndicator{Name: "Is a banana", LifeExpectancyImpact: .1},
		HealthIndicator{Name: "Eats bananas", LifeExpectancyImpact: -.2},
	}

	var expectedAgeAtDeath = pandaMaxLife - uint(.1 * float32(pandaMaxLife-testHealthData.Age))
	var expectedYearsOfMortalityRemaining = expectedAgeAtDeath - testHealthData.Age

	validateExpectedResults(testHealthData, expectedAgeAtDeath, expectedYearsOfMortalityRemaining, t)
}

func TestAnalyzePandaHealthDataWithDeadPanda(t *testing.T) {
        var testHealthData PandaHealthData
        testHealthData.Name = "TestPanda"
        testHealthData.Status = "dead"
        testHealthData.Age = 17
        testHealthData.HealthIndicators = []HealthIndicator {
                HealthIndicator{Name: "Doesn't much matter now", LifeExpectancyImpact: .1},
                HealthIndicator{Name: "Does it?", LifeExpectancyImpact: -.2},
        }

        var expectedAgeAtDeath = testHealthData.Age
        var expectedYearsOfMortalityRemaining = uint(0)

        validateExpectedResults(testHealthData, expectedAgeAtDeath, expectedYearsOfMortalityRemaining, t)
}


func validateExpectedResults(testHealthData PandaHealthData, expectedAgeAtDeath uint, expectedYearsOfMortalityRemaining uint, t *testing.T) {
	pandaHealthAnalysis := analyzePandaHealthData(testHealthData)

	if pandaHealthAnalysis.Name != testHealthData.Name {
                t.Errorf("Failed to copy name to analysis object")
        }

        if pandaHealthAnalysis.Status != testHealthData.Status {
                t.Errorf("Failed to copy status to analysis object")
        }

	if pandaHealthAnalysis.Age != testHealthData.Age {
                t.Errorf("Failed to copy age to analysis object")
        }

	if pandaHealthAnalysis.HealthIndicators == nil {
                t.Errorf("Failed to copy health indicators to analysis object")
        }

	if pandaHealthAnalysis.ExpectedAgeAtDeath != expectedAgeAtDeath {
                t.Errorf("Failed to calculate correct expected age at death, got %d but expected %d", pandaHealthAnalysis.ExpectedAgeAtDeath, expectedAgeAtDeath)
        }

	if pandaHealthAnalysis.ExpectedYearsOfMortalityRemaining != expectedYearsOfMortalityRemaining {
                t.Errorf("Failed to calculate correct expected years of life remaining, got %d but expected %d", pandaHealthAnalysis.ExpectedAgeAtDeath, expectedYearsOfMortalityRemaining)
        }

}
