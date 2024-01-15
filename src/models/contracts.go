package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

// Must implement this
type IContract interface {
	ConditionMet() bool
	Consequence()
	Reward()
}

// need to have ID
type BasicContract struct {
	gorm.Model
	// IContract // TODO can't have interface here without some special gorm deserialization logic
	UserID            uint                 `gorm:"embedded"`
	EvaluationHistory []ContractEvaluation `gorm:"foreignKey:ContractId"`
}

type ContractEvaluation struct {
	gorm.Model
	ContractId     uint // foreign key this
	EvaluationTime time.Time
	ThresholdMet   bool
	// TODO add specifics?? numbers by how much goal was missed?
}

func (c *BasicContract) ConditionMet() bool {
	return false
}

func (c *BasicContract) Consequence() {
	fmt.Printf("GOAL NOT MET for contract %d! Doing consequence", c.ID)
	return
}

func (c *BasicContract) Reward() {
	return
}

// two basic types of contracts: SingleOccurence and Recurring
type SingleOccurenceContract struct {
	BasicContract
	EvaluationTime time.Time
}

type RecurringContract struct {
	BasicContract          `gorm:"embedded"`
	EvaluationSchedule     string        // cron syntax, how often to evaluate
	EvaluationLookback     time.Duration // how long to look back, e.g. 1w
	EvaluationFunctionName string        // The name of the function to be evaluated in scheduleFunctions map
}

func (c *RecurringContract) ConditionMet() bool {
	return false
}
