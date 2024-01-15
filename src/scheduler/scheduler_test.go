package scheduler

import (
	"testing"
	"time"

	"github.com/devdetour/ulysses/models"
	"github.com/devdetour/ulysses/test"
	"github.com/stretchr/testify/assert"
)

func TestScheduleCron(t *testing.T) {
	test.Setup()

	SetupScheduler()

	// Invalid cron schedule returns err
	_, err := ScheduleCron("", func(c *models.RecurringStravaContract) {}, &models.RecurringStravaContract{}, false)
	// assert.Nil(t, jobPtr)
	assert.NotNil(t, err)

	// Setup test function
	// done := make(chan struct{})

	val := 0
	testFunc := func(c *models.RecurringStravaContract) {
		val++
	}

	// Valid cron schedule works
	job, err := ScheduleCron("* * * * *", testFunc, &models.RecurringStravaContract{}, true)

	assert.NotNil(t, job)
	assert.Nil(t, err)

	time.Sleep(1 * time.Second)

	assert.Equal(t, 1, val)
}

func TestMapNameToFunction(t *testing.T) {
	SetupNameToFunction()
	// default works
	f, err := MapNameToFunction(DEFAULT_FUNCTION_NAME)
	assert.NotNil(t, f)
	assert.Nil(t, err)

	// nil doesn't
	f, err = MapNameToFunction("asdf")
	assert.Nil(t, f)
	assert.NotNil(t, err)
}

func TestScheduleStartup(t *testing.T) {
	test.Setup()
	SetupScheduler()
	time.Sleep(10 * time.Second)
}
