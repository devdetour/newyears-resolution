package connector

import (
	"fmt"
	"net/http"

	"github.com/valyala/fasthttp"
)

const PORT int = 5000
const CONSEQUENCE_ENDPOINT_FORMAT string = "http://localhost:%d/start"
const OK_ENDPOINT_FORMAT string = "http://localhost:%d/goal_met"

// send a request to FLASK server for the user with the punishment they gotta do lol
func PunishUser(userid uint) bool {
	url := fmt.Sprintf(CONSEQUENCE_ENDPOINT_FORMAT, PORT)

	// Create a fasthttp request object
	req := fasthttp.AcquireRequest()

	// Set the request URI
	req.SetRequestURI(url)

	// Create a fasthttp response object
	resp := fasthttp.AcquireResponse()

	// Perform the GET request
	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Println("Error:", err)
		return false
	}

	// Print the response status code and body
	fmt.Printf("Status Code: %d\n", resp.StatusCode())
	fmt.Printf("Response Body: %s\n", resp.Body())

	// Release resources
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
	return resp.StatusCode() == http.StatusOK
}

func ReportOK(userid uint) bool {
	url := fmt.Sprintf(OK_ENDPOINT_FORMAT, PORT)

	// Create a fasthttp request object
	req := fasthttp.AcquireRequest()

	// Set the request URI
	req.SetRequestURI(url)

	// Create a fasthttp response object
	resp := fasthttp.AcquireResponse()

	// Perform the GET request
	if err := fasthttp.Do(req, resp); err != nil {
		fmt.Println("Error:", err)
		return false
	}

	// Print the response status code and body
	fmt.Printf("Status Code: %d\n", resp.StatusCode())
	fmt.Printf("Response Body: %s\n", resp.Body())

	// Release resources
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
	return resp.StatusCode() == http.StatusOK
}
