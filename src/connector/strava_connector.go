package connector

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/devdetour/ulysses/auth"
	"github.com/devdetour/ulysses/database"
	"github.com/devdetour/ulysses/models"
	"github.com/devdetour/ulysses/util"
	"github.com/valyala/fasthttp"
)

const (
	DISTANCE_GOAL int = iota
	TIME_GOAL
)

const ACTIVITY_URL = "https://www.strava.com/api/v3/athlete/activities"

// Call Strava API for given contract
func GetStravaDataForUserId(userId uint) ([]models.Activity, error) {

	// Get token from DB for user
	token, err := database.GetExternalAuthTokenForUser(userId)

	if err != nil {
		return nil, err
	}

	// token := ExternalAuthToken{}

	if len(token.Text) == 0 {
		return nil, fmt.Errorf("Token empty!")
	}

	// Check token scopes
	if !strings.Contains(token.Scope, "activity:read") {
		return nil, fmt.Errorf("Token has insufficient scopes!")
	}

	// Check token is not expired. Refresh if it is
	if token.Expires.Before(time.Now()) {
		token, err = auth.StravaTokenRefresh(*token)
		if err != nil {
			return nil, fmt.Errorf("Failed to refresh token!")
		}
	}

	// Get data
	data, err := getStravaData(token.Text)

	if err != nil {
		return nil, err
	}

	return data, nil
}

// Just need token right? what if expired?
func getStravaData(token string) ([]models.Activity, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(ACTIVITY_URL)
	req.Header.SetMethod("GET")

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	// Access the response body and headers
	body := resp.Body()
	statusCode := resp.StatusCode()

	fmt.Printf("Strava Status Code: %d\n", statusCode)
	// fmt.Printf("Response Body: %s\n", body)

	if !util.Is2XXStatus(statusCode) {
		return nil, fmt.Errorf("Got status code %d", statusCode)
	}

	var activities []models.Activity
	if err := json.Unmarshal(body, &activities); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}

	return activities, nil
}
