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

	response, err := utils.HttpGetRequest("https://api.tp-staging.com/v1/private/me", tokenData.AccessToken)
	if err != nil {
		fmt.Println("Error making API request:", err)
		// Handle the error appropriately
	}

	userInfo, err := utils.ProcessUserResponse(response)
	if err != nil {
		fmt.Println("Error parsing user info:", err)
	}

	buidUrl := fmt.Sprintf("https://api.tp-staging.com/v1/private/business-users/%s/business-units", userInfo.BusinessUser.ID)
	businessUnitsResponse, err := utils.HttpGetRequest(buidUrl, tokenData.AccessToken)

	fmt.Println("User name:", userInfo.BusinessUser.Name)
	fmt.Println("User email:", userInfo.BusinessUser.Email)

	if err != nil {
		fmt.Println("Error making business units API request:", err)
		http.Error(w, "Error fetching business units", http.StatusInternalServerError)
		return
	}
	businessUnitsInfo, err := utils.ProcessBusinessUnitsResponse(businessUnitsResponse)
	if err != nil {
		fmt.Println("Error parsing business units info:", err)
		http.Error(w, "Error parsing business units info", http.StatusInternalServerError)
		return
	}
	fmt.Println("Business Units:")

	for _, bu := range businessUnitsInfo.BusinessUnits {
		buidInfoUrl := fmt.Sprintf("https://api.tp-staging.com/v1/private/business-units/%s", bu.ID)
		time.Sleep(500 * time.Millisecond) // Add delay to avoid throttling
		buInfoResponse, err := utils.HttpGetRequest(buidInfoUrl, tokenData.AccessToken)
		if err != nil {
			fmt.Println("Error making business unit info API request:", err)
			http.Error(w, "Error fetching business unit info", http.StatusInternalServerError)
			return
		}

		businessUnitDetails, err := utils.ProcessBusinessUnitDetails(buInfoResponse)
		if err != nil {
			fmt.Println("Error parsing business unit details:", err)
			http.Error(w, "Error parsing business unit details", http.StatusInternalServerError)
			return
		}
		fmt.Printf("Business Unit ID: %s\n", businessUnitDetails.ID)
		fmt.Printf("Business Unit Display Name: %s\n", businessUnitDetails.DisplayName)
		state.AppendToBuids(businessUnitDetails.DisplayName)
	}
}
