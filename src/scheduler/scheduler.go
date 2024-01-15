package scheduler

import (
	"fmt"
	"log"

	"github.com/devdetour/ulysses/database"
	"github.com/devdetour/ulysses/models"
	"github.com/go-co-op/gocron/v2"
)

var Scheduler gocron.Scheduler

func SetupScheduler() {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatal(err)
	}
	Scheduler = s
	// something something load all the jobs
	db := database.DB

	// find all schedules in db
	var contracts []models.RecurringStravaContract
	db.Find(&contracts)

	for _, c := range contracts {
		// Read what function each contract SHOULD have, and use that function if it exists. Error if not
		evalauteFunc, err := MapNameToFunction(c.EvaluationFunctionName)
		if err != nil {
			log.Fatalf("Failed to get function for contract! Contract ID: %d, function name: %s", c.ID, c.EvaluationFunctionName)
		}
		// copy c so pointer isn't updated by next run of the loop
		copy := c
		ScheduleCron(c.RecurringContract.EvaluationSchedule, evalauteFunc, &copy, false) // TODO un-true this for prod
	}

	Scheduler.Start()
}

func RestartScheduler() {
	Scheduler.Shutdown()
	SetupScheduler()
}

func StartScheduler() {
	return
}

// Schedule a job with CRON syntax
func ScheduleCron(cronSchedule string, scheduleFunc func(c *models.RecurringStravaContract), c *models.RecurringStravaContract, startImmediately bool) (gocron.Job, error) {
	// create func to wrap contract.Evaluate, then schedule
	if scheduleFunc == nil {
		return nil, fmt.Errorf("Nil schedulefunc")
	}

	options := []gocron.JobOption{}

	if startImmediately {
		options = append(options, gocron.JobOption(gocron.WithStartImmediately()))
	}

	job, err := Scheduler.NewJob(
		gocron.CronJob(cronSchedule, false), // TODO support with seconds, check cron schedule to see how many tokens
		gocron.NewTask(scheduleFunc, c),
		options...)
	return job, err

	// job, err := Scheduler.NewJob(
	// 	gocron.CronJob(cronSchedule, false), // TODO support with seconds, check cron schedule to see how many tokens
	// 	gocron.NewTask(scheduleFunc))
	// return job, err

	// Cron(cronSchedule).StartImmediately().Do(scheduleFunc)
}
