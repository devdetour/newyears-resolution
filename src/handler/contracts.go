package handler

import (
	"fmt"
	"time"

	"github.com/devdetour/ulysses/database"
	"github.com/devdetour/ulysses/models"
	"github.com/devdetour/ulysses/scheduler"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type ContractInput struct {
	Type         string  `json:"type"` // recurring or one-shot
	GoalCategory string  `json:"goalCategory"`
	GoalType     int     `json:"goalType"`
	Goal         float64 `json:"goal"`
	Lookback     int     `json:"lookback"` // TODO time interval??
	LookbackUnit string  `json:"lookbackUnit"`
	Schedule     string  `json:"schedule"`
}

// Returns true if all fields are valid
func validateContract(i ContractInput) bool {
	return i.GoalCategory == "strava" &&
		(i.GoalType >= 0 && i.GoalType < 2) &&
		(i.Type == "recurring") && // can be one shot to... TODO
		i.Lookback > 0 &&
		len(i.Schedule) > 0 // TODO validate real cron
}

func CreateContract(c *fiber.Ctx) error {
	// {"type":"recurring","schedule":"* * * * *","goalCategory":"strava","goalType":1,"goal":100,"lookback":100}
	db := database.DB

	// Make sure token valid
	// TODO make all this a helper methodddd
	token := c.Locals("user").(*jwt.Token)

	var userId uint
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userId = uint(claims["user_id"].(float64))
		if userId == 0 {
			return fmt.Errorf("Must have username")
		}
	}

	fmt.Print(userId) // Create contract for user

	input := new(ContractInput)

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error creating contract", "data": err})
	}

	if !validateContract(*input) {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error creating contract", "data": "Bad input!"})
		return nil
	}

	// make sure schedule is valid
	if len(input.Schedule) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Schedule required!"})
	}

	// make sure function name valid & exists
	if len(input.GoalCategory) == 0 {
		_, err := scheduler.MapNameToFunction(input.GoalCategory) // TODO can probably combine onto prev. line
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Function name invalid!", "data": err})
		}
	}

	// make sure lookback unit is valid
	if input.LookbackUnit != "hours" { // TODO not just hours
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Lookback unit invalid!", "data": fmt.Errorf("Make sure it is hours")})
	}

	lookbackUnit := time.Hour

	contract := models.RecurringStravaContract{
		RecurringContract: models.RecurringContract{
			BasicContract: models.BasicContract{
				UserID: 1,
			},
			EvaluationSchedule:     input.Schedule,
			EvaluationFunctionName: input.GoalCategory,
			EvaluationLookback:     time.Duration(input.Lookback) * lookbackUnit,
		},
		StravaContract: models.StravaContract{
			Goal: models.StravaGoal{
				GoalThreshold: input.Goal,
				GoalType:      input.GoalType,
			},
		},
	}

	tx := db.Create(&contract)
	fmt.Print(tx)

	// If creation succeeded, also schedule contract

	// Try to get function from store
	evalFn, err := scheduler.MapNameToFunction(contract.EvaluationFunctionName)

	if err != nil {
		return fmt.Errorf("Failed to get eval function! %v", err)
	}

	scheduler.ScheduleCron(contract.EvaluationSchedule, evalFn, &contract, false)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Contract created!"})

}

func DeleteContract(c *fiber.Ctx) error {
	// {"type":"recurring","schedule":"* * * * *","goalCategory":"strava","goalType":1,"goal":100,"lookback":100}
	db := database.DB

	userId, err := getJwtUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error deleting contract", "data": "Couldn't get userId!"})
	}

	fmt.Print(userId)

	type DeleteContractInput struct {
		ContractId uint `json:"contractId"`
	}

	input := new(DeleteContractInput)

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error deleting contract", "data": err})
	}

	// make sure contract exists
	var contract models.RecurringStravaContract
	result := db.Find(&contract, "User_Id = ? AND ID = ?", userId, input.ContractId)
	if result.RowsAffected < 1 || result.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error deleting contract", "data": fmt.Sprintf("Contract with id %d not found!", input.ContractId)})
	}

	// Make sure contract user in DB matches user from request
	if contract.RecurringContract.BasicContract.UserID != userId {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "error", "message": "Error deleting contract", "data": fmt.Sprintf("User does not own contract with ID %d!", input.ContractId)})
	}

	result = db.Delete(&contract, "ID = ?", input.ContractId)

	// Unschedule job for contract if deleted
	// TODO better way of doing this!
	scheduler.RestartScheduler()
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": fmt.Sprintf("Deleted contract %d", input.ContractId)})

}

func GetContractsForUser(c *fiber.Ctx) error {
	db := database.DB

	// Make sure token valid
	// TODO make all this a helper methodddd
	token := c.Locals("user").(*jwt.Token)

	var userId uint
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		userId = uint(claims["user_id"].(float64))
		if userId == 0 {
			fmt.Errorf("Must have username")
		}
	}

	fmt.Print(userId) // Create contract for user

	var contractList []models.RecurringStravaContract // TODO eventually change this to just regular contracts
	result := db.Find(&contractList, "User_Id = ?", userId)

	fmt.Print(result)

	c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok", "message": "Got contracts for user", "data": contractList})
	return nil
}

func GetEvaluationHistory(c *fiber.Ctx) error {
	db := database.DB
	userId, err := getJwtUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "JWT token invalid", "data": ""})
	}

	// Get contracts
	var contractIds []uint
	var allContracts []models.RecurringStravaContract                                                      // TODO eventually change this to just regular contracts
	result := db.Model(&models.RecurringStravaContract{}).Where("User_Id = ?", userId).Find(&allContracts) //.Pluck("ID", &contractIds)
	// result = db.Model(&models.RecurringStravaContract{}).Pluck("ID", &contractIds) // WORKS

	result = db.Model(&models.RecurringStravaContract{}).Pluck("ID", &contractIds).Where("User_Id = ?", userId)

	var contractHistory []models.ContractEvaluation
	result = db.Where("contract_id IN ?", contractIds).Find(&contractHistory)
	// Get contract history

	fmt.Print(result)
	c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "ok", "message": "Got contracts for user", "data": contractHistory})
	return nil
}
