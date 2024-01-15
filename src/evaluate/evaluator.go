package evaluate

import (
	"github.com/devdetour/ulysses/models"
)

// why did i do this?????
func ConditionMet(c *models.BasicContract) bool {
	return c.ConditionMet()
}
