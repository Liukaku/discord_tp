package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/liukaku/discord-tp/cmd/server/types"
)

func HttpGetRequest(url string, bearer string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating GET request: %w", err)
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making GET request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}

	return string(body), nil
}

func HttpPostRequest(url string, bearer string, body string) (string, error) {
	postBody := strings.NewReader(body)
	req, err := http.NewRequest("POST", url, postBody)
	if err != nil {
		return "", fmt.Errorf("error creating POST request: %w", err)
	}
	if bearer != "" {
		req.Header.Set("Authorization", "Bearer "+bearer)
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making POST request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 response code: %d", resp.StatusCode)
	}
	respBod, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %w", err)
	}
	return string(respBod), nil
}

func SendReviewReply(reviewId string, userId string, replyMessage string, bearer string) error {
	url := fmt.Sprintf("https://api.tp-staging.com/v1/private/reviews/%s/reply", reviewId)
	body := fmt.Sprintf(`{"authorBusinessUnserId": "%s", "message": "%s"}`, userId, replyMessage)

	resp, err := HttpPostRequest(url, bearer, body)
	if err != nil {
		return fmt.Errorf("error sending review reply: %w", err)
	}

	fmt.Println("Review reply response:", resp)
	return nil
}

func ProcessUserResponse(responseBody string) (*types.UserResponse, error) {
	var userResponse types.UserResponse
	err := json.Unmarshal([]byte(responseBody), &userResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing user response: %w", err)
	}

	return &userResponse, nil
}

func ProcessBusinessUnitsResponse(responseBody string) (*types.BusinessUnitsResponse, error) {
	var businessUnitsResponse types.BusinessUnitsResponse
	err := json.Unmarshal([]byte(responseBody), &businessUnitsResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing business units response: %w", err)
	}

	return &businessUnitsResponse, nil
}

func ProcessBusinessUnitDetails(responseBody string) (*types.BusinessUnitDetails, error) {
	var businessUnitDetails types.BusinessUnitDetails
	err := json.Unmarshal([]byte(responseBody), &businessUnitDetails)
	if err != nil {
		return nil, fmt.Errorf("error parsing business unit details: %w", err)
	}

	return &businessUnitDetails, nil
}

func ConvertTrustpilotApiUrlToPublic(apiUrl string) string {
	// First, extract the review ID
	re := regexp.MustCompile(`https://api\.tp-staging\.com/v1/reviews/([a-f0-9]+)`)
	matches := re.FindStringSubmatch(apiUrl)

	if len(matches) < 2 {
		fmt.Println("Warning: Could not extract review ID from URL:", apiUrl)
		return apiUrl // Return original if no match
	}

	// reviewId := matches[1]

	// Create the public URL
	publicUrl := strings.Replace(apiUrl, "api.tp-staging.com/v1", "www.tp-staging.com", 1)

	return publicUrl
}
