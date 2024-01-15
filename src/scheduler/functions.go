package scheduler

import (
	"fmt"
	"time"

	"github.com/devdetour/ulysses/connector"
	"github.com/devdetour/ulysses/database"
	"github.com/devdetour/ulysses/models"
)

const (
	DEFAULT_FUNCTION_NAME           string = "default"
	STRAVA_ACTIVITIES_FUNCTION_NAME string = "strava"
)

var scheduleFunctions map[string]func(c *models.RecurringStravaContract)

func MapNameToFunction(functionName string) (func(c *models.RecurringStravaContract), error) {
	f, ok := scheduleFunctions[functionName]
	if !ok {
		return nil, fmt.Errorf("Function %s not found!", functionName)
	}
	return f, nil
}

func SetupNameToFunction() {
	// Map functions to names
	scheduleFunctions = make(map[string]func(c *models.RecurringStravaContract))
	scheduleFunctions[DEFAULT_FUNCTION_NAME] = defaultFn
	scheduleFunctions[STRAVA_ACTIVITIES_FUNCTION_NAME] = CheckStravaActivities
}

func defaultFn(c *models.RecurringStravaContract) {
	// Set the time zone to PST
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	// Get the current time in the specified time zone
	currentTime := time.Now().In(location)

	// Print the current time
	fmt.Println("RAN! Current time in PST:", currentTime.Format("2006-01-02 15:04:05 MST"))
}

// Function that takes in a contract and returns an argumentless function to get the Strava activity data for that contract's user
func CheckStravaActivities(c *models.RecurringStravaContract) {
	defaultFn(c)
	// Get activities
	activities, err := connector.GetStravaDataForUserId(c.BasicContract.UserID)
	if err != nil {
		fmt.Errorf("Failed to get strava data for user %d: %v", c.BasicContract.UserID, err)
		return
	}

	if len(activities) == 0 {
		fmt.Errorf("No activities found for user %d!", c.UserID)
		return
	}

	// If we got back activities, check if goal met
	goalMet := c.GoalMet(time.Now(), activities)

	CreateEvaluationRecord(&c.BasicContract, goalMet)

	if !goalMet {
		connector.PunishUser(c.UserID)
		// c.Consequence()
	} else {
		connector.ReportOK(c.UserID)
	}

	fmt.Printf("Goal met for contract ID %d: %t\n", c.ID, goalMet)
}

// function that wraps a call to strava API for renewing token
func RenewStravaToken() func() {
	// Set the time zone to PST
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return nil
	}

	// Get the current time in the specified time zone
	currentTime := time.Now().In(location)

	// Print the current time
	fmt.Println("RAN! Current time in PST:", currentTime.Format("2006-01-02 15:04:05 MST"))

	return nil
}

func CreateEvaluationRecord(c *models.BasicContract, goalMet bool) {
	// Create history
	evaluationRecord := models.ContractEvaluation{
		ContractId:     c.ID,
		EvaluationTime: time.Now(),
		ThresholdMet:   goalMet,
	}
	db := database.DB
	db.Create(&evaluationRecord)
}
