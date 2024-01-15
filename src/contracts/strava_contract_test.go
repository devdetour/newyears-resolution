package contracts

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/devdetour/ulysses/database"
	"github.com/devdetour/ulysses/models"
	"github.com/devdetour/ulysses/test"
	"github.com/stretchr/testify/assert"
)

func Test_GetStravaData(t *testing.T) {
	test.Setup()

	// TODO I don't like this syntax really... better way?
	c := models.RecurringStravaContract{
		RecurringContract: models.RecurringContract{
			BasicContract: models.BasicContract{
				UserID: 1,
			},
			EvaluationSchedule: "10",
			EvaluationLookback: time.Hour * 24,
		},
		StravaContract: models.StravaContract{
			Goal: models.StravaGoal{
				GoalType:      models.DISTANCE_GOAL,
				GoalThreshold: 100.2,
			},
		},
	}

	db := database.DB
	err := db.Create(&c).Error
	assert.Nil(t, err)

	// TODO test for getting multiple contracts
	var result models.RecurringStravaContract
	db.Find(&result, "user_id = ?", 1)

	assert.NotNil(t, result)
	assert.Equal(t, result.StravaContract.Goal.GoalType, c.StravaContract.Goal.GoalType)
	assert.Equal(t, c.StravaContract.Goal.GoalThreshold, result.StravaContract.Goal.GoalThreshold)

	fmt.Print(c)
	test.TearDown()
}

func Test_GoalMet(t *testing.T) {
	activities := make([]models.Activity, 1)
	if err := json.Unmarshal([]byte(testData), &activities); err != nil {
		fmt.Println("Error parsing JSON:", err)
		assert.True(t, false)
		return
	}

	assert.Len(t, activities, 8)

	c := models.RecurringStravaContract{
		RecurringContract: models.RecurringContract{
			BasicContract: models.BasicContract{
				UserID: 1,
			},
			EvaluationSchedule: "10",
			EvaluationLookback: time.Hour * 24,
		},
		StravaContract: models.StravaContract{
			Goal: models.StravaGoal{
				GoalType:      models.DISTANCE_GOAL,
				GoalThreshold: 100.2,
			},
		},
	}

	tests := []struct {
		EvalTime  time.Time
		Threshold float64
		GoalType  int
		Expected  bool
	}{
		{
			// Distance goal met based on activities within lookback window
			time.Date(2023, time.March, 20, 1, 23, 39, 0, time.UTC),
			100.0,
			models.DISTANCE_GOAL,
			true,
		},
		{
			// Distance goal NOT met (activities too far out of time window)
			time.Date(2023, time.June, 20, 1, 23, 39, 0, time.UTC),
			100.0,
			models.DISTANCE_GOAL,
			false,
		},
		{
			// Distance goal NOT met (activity distance too small)
			time.Date(2023, time.March, 20, 1, 23, 39, 0, time.UTC),
			100000.0,
			models.DISTANCE_GOAL,
			false,
		},
		{
			// Time goal met
			time.Date(2023, time.March, 20, 1, 23, 39, 0, time.UTC),
			100.0,
			models.TIME_GOAL,
			true,
		},
		{
			// Time goal NOT met
			time.Date(2023, time.March, 20, 1, 23, 39, 0, time.UTC),
			10000.0,
			models.TIME_GOAL,
			false,
		},
	}

	for _, test := range tests {
		c.Goal.GoalThreshold = test.Threshold
		c.Goal.GoalType = test.GoalType
		assert.Equal(t, test.Expected, c.GoalMet(test.EvalTime, activities))
	}
}

const testData = `
 <<< my own strava activity data redacted, fill in your own strava activites (json list format) >>>
`
