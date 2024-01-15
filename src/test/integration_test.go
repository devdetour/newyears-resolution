package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/devdetour/ulysses/server"
	"github.com/stretchr/testify/assert"
)

var USER_LOGIN_DATA map[string]string = map[string]string{
	"username": "test",
	"email":    "test@test.test",
	"password": "123456",
}

type LoginResponse struct {
	Data    string `json:"data"`
	Message string `json:"message"`
	Status  string `json:"status"`
}

func SendCreateContractRequest(jwt string) int {
	// yeet
	url := "http://localhost:3000/api/contracts/create"

	type CreateContractPayload struct {
		Type         string `json:"type"`
		Schedule     string `json:"schedule"`
		GoalCategory string `json:"goalCategory"`
		GoalType     int    `json:"goalType"`
		Goal         int    `json:"goal"`
		Lookback     int    `json:"lookback"`
	}

	data := CreateContractPayload{
		Type:         "recurring",
		Schedule:     "* * * * *",
		GoalCategory: "strava",
		GoalType:     0,
		Goal:         1000000,
		Lookback:     10,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return -1
	}

	// Create a new HTTP request with the POST method
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return -1
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Set the Authorization header with the bearer token
	req.Header.Set("Authorization", "Bearer "+jwt)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return -1
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	// Handle response body if needed
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return -1
	}
	fmt.Println("Response Body:", string(body))
	return resp.StatusCode
}

func SendCreateUserRequest() {
	// yeet
	url := "http://localhost:3000/auth/register"
	data := USER_LOGIN_DATA

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	// Handle response body if needed
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	fmt.Println("Error reading response body:", err)
	// 	return
	// }
	// fmt.Println("Response Body:", string(body))
}

func LoginUser(t *testing.T) (string, error) {
	// yeet
	url := "http://localhost:3000/auth/login"
	data := map[string]string{
		"identity": "test",
		"password": "123456",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return "", err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return "", err
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	assert.Equal(t, 200, resp.StatusCode)

	// Handle response body if needed
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return "", err
	}

	var response LoginResponse
	err = json.Unmarshal([]byte(string(body)), &response)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", err
	}

	fmt.Println("Response Body:", string(body))

	return response.Data, nil
}

func SendGetContractRequest(jwt string) {
	url := "http://localhost:3000/api/contracts/get"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the content type header
	req.Header.Set("Content-Type", "application/json")

	// Set the Authorization header with the bearer token
	req.Header.Set("Authorization", "Bearer "+jwt)
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)

	// Handle response body if needed
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println("Response Body:", string(body))
	return
}

func Test_Integration_ContractEvaluation(t *testing.T) {
	// Clean DB
	Setup()

	// Start server
	go server.Run()

	time.Sleep(1 * time.Second)

	// Create contract to test, with goal
	SendCreateUserRequest()
	jwt, err := LoginUser(t)

	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(jwt), 0)

	// Fail test if error
	if err != nil {
		assert.True(t, false)
	}

	// Create contract using JWT
	status := SendCreateContractRequest(jwt)
	assert.Equal(t, 200, status)

	// Get contracts, assert new contract was created
	SendGetContractRequest(jwt)

	// TODO this
	// contracts := SendGetContractRequest(jwt)
	// assert.NotNil(t, contracts)
	// assert.Equal(t, len(contracts), 1)

	TearDown()
}
