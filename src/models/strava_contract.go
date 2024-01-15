package models

import (
	"fmt"
	"time"
)

const (
	DISTANCE_GOAL int = iota
	TIME_GOAL
)

const ACTIVITY_URL = "https://www.strava.com/api/v3/athlete/activities"

type StravaGoal struct {
	GoalType      int
	GoalThreshold float64
}

type StravaContract struct {
	Goal StravaGoal `gorm:"embedded"`
}

type RecurringStravaContract struct {
	RecurringContract `gorm:"embedded"`
	StravaContract    `gorm:"embedded"`
}

func (c *RecurringStravaContract) GoalMet(evalTime time.Time, activities []Activity) bool {
	// Recurring contract: goal is either a workout time, or a workout distance, so track both
	totalDistance := 0.0
	totalTime := 0

	windowStart := evalTime.Add(-c.EvaluationLookback)
	// Find activities in the given evaluation window
	for _, activity := range activities {
		if activityInWindow(windowStart, activity) {
			totalDistance += activity.Distance
			totalTime += activity.MovingTime
		}
	}

	if c.Goal.GoalType == DISTANCE_GOAL {
		return totalDistance >= c.Goal.GoalThreshold
	}

	if c.Goal.GoalType == TIME_GOAL {
		return float64(totalTime)/60 >= c.Goal.GoalThreshold
	}

	return false
}

func (c *RecurringStravaContract) Consequence() {
	// Do a thing!

	fmt.Print("Criteria not met! consequence happening")
	return
}

func activityInWindow(startTime time.Time, activity Activity) bool {
	start := activity.StartDate
	end := activity.StartDate.Add(time.Duration(activity.MovingTime) * time.Second)
	// Check if start or end time are within window
	return startTime.Before(start) || startTime.Before(end)
}
