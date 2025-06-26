package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/liukaku/discord-tp/cmd/server/handlers/utils"
	"github.com/liukaku/discord-tp/cmd/server/types"
)

func StoreTokenHandler(w http.ResponseWriter, r *http.Request, state types.SharedState) {
	var tokenData struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   string `json:"expires_in"`
	}

	// Parse the JSON body
	err := json.NewDecoder(r.Body).Decode(&tokenData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	fmt.Println("Received access token:", tokenData.AccessToken)
	fmt.Println("Token type:", tokenData.TokenType)
	fmt.Println("Expires in:", tokenData.ExpiresIn)

	userInfo, err := getAboutMe(tokenData.AccessToken)

	if err != nil {
		http.Error(w, "Error fetching user info", http.StatusInternalServerError)
		return
	}

	businessUnitsInfo, err := getBusinessUnitsInfo(userInfo.BusinessUser.ID, tokenData.AccessToken)

	if err != nil {
		http.Error(w, "Error fetching business units info", http.StatusInternalServerError)
		return
	}

	for _, bu := range businessUnitsInfo.BusinessUnits {
		businessUnitDetails, err := getSingleBusinessUnitInfo(bu.ID, tokenData.AccessToken)

		if err != nil {
			http.Error(w, "Error fetching single business unit info", http.StatusInternalServerError)
			return
		}

		state.AppendToBuids(businessUnitDetails.DisplayName)
	}
}

func getAboutMe(accessToken string) (*types.UserResponse, error) {
	response, err := utils.HttpGetRequest("https://api.tp-staging.com/v1/private/me", accessToken)
	if err != nil {
		fmt.Println("Error making API request:", err)
		return nil, err
	}

	userInfo, err := utils.ProcessUserResponse(response)
	if err != nil {
		fmt.Println("Error parsing user info:", err)
		return nil, err
	}

	return userInfo, nil
}

func getBusinessUnitsInfo(businessUserId string, accessToken string) (*types.BusinessUnitsResponse, error) {
	buidUrl := fmt.Sprintf("https://api.tp-staging.com/v1/private/business-users/%s/business-units", businessUserId)
	businessUnitsResponse, err := utils.HttpGetRequest(buidUrl, accessToken)

	if err != nil {
		fmt.Println("Error making business units API request:", err)
		return nil, err
	}
	businessUnitsInfo, err := utils.ProcessBusinessUnitsResponse(businessUnitsResponse)
	if err != nil {
		fmt.Println("Error parsing business units info:", err)
		return nil, err
	}
	fmt.Println("Business Units:")

	return businessUnitsInfo, nil
}

func getSingleBusinessUnitInfo(businessUnitId string, accessToken string) (*types.BusinessUnitDetails, error) {
	buidInfoUrl := fmt.Sprintf("https://api.tp-staging.com/v1/private/business-units/%s", businessUnitId)
	time.Sleep(500 * time.Millisecond) // Add delay to avoid throttling
	buInfoResponse, err := utils.HttpGetRequest(buidInfoUrl, accessToken)
	if err != nil {
		fmt.Println("Error making business unit info API request:", err)
		return nil, err
	}

	businessUnitDetails, err := utils.ProcessBusinessUnitDetails(buInfoResponse)
	if err != nil {
		fmt.Println("Error parsing business unit details:", err)
		return nil, err
	}
	fmt.Printf("Business Unit ID: %s\n", businessUnitDetails.ID)
	fmt.Printf("Business Unit Display Name: %s\n", businessUnitDetails.DisplayName)
	return businessUnitDetails, nil
}
