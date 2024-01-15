package auth

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/devdetour/ulysses/config"
	"github.com/devdetour/ulysses/models"
	"github.com/valyala/fasthttp"
)

var (
	CLIENT_SECRET string = config.Config("STRAVA_CLIENT_SECRET")
	CLIENT_ID     string = config.Config("STRAVA_CLIENT_ID")
)

type Athlete struct {
	ID            int         `json:"id"`
	Username      string      `json:"username"`
	ResourceState int         `json:"resource_state"`
	Firstname     string      `json:"firstname"`
	Lastname      string      `json:"lastname"`
	Bio           string      `json:"bio"`
	City          string      `json:"city"`
	State         string      `json:"state"`
	Country       string      `json:"country"`
	Sex           string      `json:"sex"`
	Premium       bool        `json:"premium"`
	Summit        bool        `json:"summit"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	BadgeTypeID   int         `json:"badge_type_id"`
	Weight        float64     `json:"weight"`
	ProfileMedium string      `json:"profile_medium"`
	Profile       string      `json:"profile"`
	Friend        interface{} `json:"friend"`
	Follower      interface{} `json:"follower"`
}

type TokenResponse struct {
	TokenType    string  `json:"token_type"`
	ExpiresAt    int     `json:"expires_at"`
	ExpiresIn    int     `json:"expires_in"`
	RefreshToken string  `json:"refresh_token"`
	AccessToken  string  `json:"access_token"`
	Athlete      Athlete `json:"athlete"`
}

// TODO maybe this should return a whole object... so we get refresh token & such.
func StravaTokenExchange(code string) (*TokenResponse, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	url := "https://www.strava.com/oauth/token"

	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	queryParams := fasthttp.Args{}
	queryParams.Add("client_id", CLIENT_ID)
	queryParams.Add("client_secret", CLIENT_SECRET)
	queryParams.Add("code", code)
	queryParams.Add("grant_type", "authorization_code")
	req.URI().SetQueryString(queryParams.String())

	err := fasthttp.Do(req, resp)
	body := resp.Body()

	fmt.Print("client_id", CLIENT_ID)
	fmt.Print("client_secret", CLIENT_SECRET)

	if err != nil || resp.StatusCode() == 400 {
		fmt.Printf("Failed to get Strava token! Status: %d\n", resp.StatusCode())
		return nil, fmt.Errorf(string(body))
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}

	// Print the response
	fmt.Printf("Status Code: %d\n", resp.StatusCode())
	fmt.Printf("Response Body: %s\n", body)
	return &tokenResponse, nil
}

// TODO probably combine with above
func StravaTokenRefresh(token models.ExternalAuthToken) (*models.ExternalAuthToken, error) {
	url := "https://www.strava.com/api/v3/oauth/token"
	clientID := CLIENT_ID
	clientSecret := CLIENT_SECRET

	refreshToken := token.RefreshToken

	// Create the request body
	requestBody := fmt.Sprintf("client_id=%s&client_secret=%s&grant_type=refresh_token&refresh_token=%s",
		clientID, clientSecret, refreshToken)

	// Create a fasthttp POST request
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBodyString(requestBody)

	// Create a fasthttp client and send the request
	client := &fasthttp.Client{}
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	// Parse the response body
	var tokenResponse TokenResponse
	body := resp.Body()
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return nil, err
	}

	// Convert to a real token
	refreshedToken := models.ExternalAuthToken{
		UserId:       token.UserId,
		Text:         tokenResponse.AccessToken,
		Source:       token.Source,
		Scope:        token.Scope,
		Expires:      time.Unix(int64(tokenResponse.ExpiresAt), 0),
		RefreshToken: tokenResponse.RefreshToken,
	}

	// Print the response
	fmt.Printf("Status Code: %d\n", resp.StatusCode())
	fmt.Printf("Response Body: %s\n", body)
	return &refreshedToken, nil
}
