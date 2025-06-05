package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// Update the shared state interface
type SharedState interface {
	GetChannelIDs() []string
	AppendToBuids(values ...string)
}

// UserResponse represents the top-level response structure
type UserResponse struct {
	BusinessUser BusinessUser `json:"businessUser"`
	Links        []Link       `json:"links"`
}

// BusinessUser represents the business user information
type BusinessUser struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Email           string    `json:"email"`
	Created         time.Time `json:"created"`
	ActivationDate  time.Time `json:"activationDate"`
	Locale          string    `json:"locale"`
	CountryID       int       `json:"countryId"`
	HasSecondFactor bool      `json:"hasSecondFactor"`
}

// Link represents a hypermedia link
type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}

// BusinessUnitsResponse represents the top-level response structure for business units
type BusinessUnitsResponse struct {
	Links         []Link         `json:"links"`
	BusinessUnits []BusinessUnit `json:"businessUnits"`
}

// BusinessUnit represents a single business unit
type BusinessUnit struct {
	Links []Link `json:"links"`
	ID    string `json:"id"`
}

// BusinessUnitDetails represents detailed information about a business unit
type BusinessUnitDetails struct {
	ID                          string                 `json:"id"`
	Country                     string                 `json:"country"`
	DisplayName                 string                 `json:"displayName"`
	HasAccountManagementConsent bool                   `json:"hasAccountManagementConsent"`
	Name                        BusinessUnitName       `json:"name"`
	Score                       Score                  `json:"score"`
	Status                      string                 `json:"status"`
	WebsiteURL                  string                 `json:"websiteUrl"`
	NumberOfReviews             ReviewCount            `json:"numberOfReviews"`
	CompanyName                 *string                `json:"companyName"`
	Description                 Description            `json:"description"`
	Address                     Address                `json:"address"`
	SocialMedia                 map[string]interface{} `json:"socialMedia"`
	Email                       *string                `json:"email"`
	Phone                       *string                `json:"phone"`
	IsClaimed                   bool                   `json:"isClaimed"`
	IsCommentsEnabled           bool                   `json:"isCommentsEnabled"`
	IsIncentivisingUsers        bool                   `json:"isIncentivisingUsers"`
	TierID                      string                 `json:"tierId"`
	FacebookPageURL             *string                `json:"facebookPageUrl"`
	IsFacebookActivated         bool                   `json:"isFacebookActivated"`
	IsSubscriber                bool                   `json:"isSubscriber"`
	FacebookPageID              int                    `json:"facebookPageId"`
	Links                       []Link                 `json:"links"`
	Warning                     string                 `json:"warning"`
	Verification                string                 `json:"verification"`
	HasSubscription             bool                   `json:"hasSubscription"`
	IsUsingPaidFeatures         bool                   `json:"isUsingPaidFeatures"`
	HideCompetitorModule        bool                   `json:"hideCompetitorModule"`
}

// BusinessUnitName represents the name information of a business unit
type BusinessUnitName struct {
	Referring   []string `json:"referring"`
	Identifying string   `json:"identifying"`
}

// Score represents the rating scores
type Score struct {
	Stars      float64 `json:"stars"`
	TrustScore float64 `json:"trustScore"`
}

// ReviewCount represents review statistics
type ReviewCount struct {
	Total                        int `json:"total"`
	UsedForTrustScoreCalculation int `json:"usedForTrustScoreCalculation"`
	OneStar                      int `json:"oneStar"`
	TwoStars                     int `json:"twoStars"`
	ThreeStars                   int `json:"threeStars"`
	FourStars                    int `json:"fourStars"`
	FiveStars                    int `json:"fiveStars"`
}

// Description represents the business description
type Description struct {
	Header *string `json:"header"`
	Text   *string `json:"text"`
}

// Address represents the business address
type Address struct {
	City        *string `json:"city"`
	State       *string `json:"state"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Postcode    *string `json:"postcode"`
	Street      *string `json:"street"`
}

func CreateHttpServer(discord *discordgo.Session, state SharedState) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}
	// channelId := os.Getenv("CHANNEL_ID")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// wg.Add(1)
		fmt.Println("received get request")
		fmt.Println("basepath URL frag:", r.URL.Fragment)
		if r.Method != http.MethodGet {
			fmt.Println("Method received: ", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// fmt.Println(r.Body)
		// // read request body
		// var resp map[string]interface{}
		// body, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	http.Error(w, "can't read body", http.StatusBadRequest)
		// 	panic(err)
		// }
		// json.Unmarshal(body, &resp)
		// fmt.Println(string(body))
		// fmt.Println(resp)
		// io.WriteString(w, fmt.Sprintf("Received request: %s \n song number: %s", r.Method, k))
		// get all guilds the bot is in
		guilds, err := discord.UserGuilds(100, "", "", true)
		if err != nil {
			fmt.Println("Error getting guilds:", err)
			http.Error(w, "Error getting guilds", http.StatusInternalServerError)
			return
		}

		channelIds := state.GetChannelIDs()
		fmt.Println("Channel IDs:")
		for i, channelId := range channelIds {
			fmt.Printf("Channel ID %d: %s\n", i, channelId)
			discord.ChannelMessageSend(channelId, "iyaaaaaaa from http")
		}

		fmt.Println("Guilds:")
		for i, guild := range guilds {
			fmt.Printf("Guild %d: %s (%s)\n", i, guild.Name, guild.ID)
		}
		// for i, guild := range guilds {
		// 	fmt.Printf("Guild %d: %s (%s)\n", i, guild.Name, guild.ID)
		// 	channels, _ := discord.GuildChannels(guild.ID)
		// 	for i, channel := range channels {
		// 		fmt.Printf("Channel %d: %s (%s)\n", i, channel.Name, channel.ID)
		// 		discord.ChannelMessageSend(channel.ID, "iyaaaaaaa")
		// 	}
		// }
		// get all channels in the guild

	})

	// /auth#access_token=tpa-f15f4140ae0a487b86be547fd76f&token_type=bearer&expires_in=360000
	http.HandleFunc("/auth", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("received auth request")
		fmt.Println("Query parameters:", r.URL.Query())
		if r.Method != http.MethodGet {
			fmt.Println("Method received: ", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// Serve the HTML page that extracts access token from URL fragment
		htmlPage := CreateAuthPage()
		w.Header().Set("Content-Type", "text/html")
		_, err := io.WriteString(w, htmlPage)
		if err != nil {
			fmt.Println("Error writing HTML page:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/store-token", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

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

		response, err := HttpGetRequest("https://api.tp-staging.com/v1/private/me", tokenData.AccessToken)
		if err != nil {
			fmt.Println("Error making API request:", err)
			// Handle the error appropriately
		}

		userInfo, err := processUserResponse(response)
		if err != nil {
			fmt.Println("Error parsing user info:", err)
		}

		buidUrl := fmt.Sprintf("https://api.tp-staging.com/v1/private/business-users/%s/business-units", userInfo.BusinessUser.ID)
		businessUnitsResponse, err := HttpGetRequest(buidUrl, tokenData.AccessToken)

		fmt.Println("User name:", userInfo.BusinessUser.Name)
		fmt.Println("User email:", userInfo.BusinessUser.Email)

		if err != nil {
			fmt.Println("Error making business units API request:", err)
			http.Error(w, "Error fetching business units", http.StatusInternalServerError)
			return
		}
		businessUnitsInfo, err := processBusinessUnitsResponse(businessUnitsResponse)
		if err != nil {
			fmt.Println("Error parsing business units info:", err)
			http.Error(w, "Error parsing business units info", http.StatusInternalServerError)
			return
		}
		fmt.Println("Business Units:")

		for _, bu := range businessUnitsInfo.BusinessUnits {
			buidInfoUrl := fmt.Sprintf("https://api.tp-staging.com/v1/private/business-units/%s", bu.ID)
			time.Sleep(500 * time.Millisecond) // Add delay to avoid throttling
			buInfoResponse, err := HttpGetRequest(buidInfoUrl, tokenData.AccessToken)
			if err != nil {
				fmt.Println("Error making business unit info API request:", err)
				http.Error(w, "Error fetching business unit info", http.StatusInternalServerError)
				return
			}

			businessUnitDetails, err := processBusinessUnitDetails(buInfoResponse)
			if err != nil {
				fmt.Println("Error parsing business unit details:", err)
				http.Error(w, "Error parsing business unit details", http.StatusInternalServerError)
				return
			}
			fmt.Printf("Business Unit ID: %s\n", businessUnitDetails.ID)
			fmt.Printf("Business Unit Display Name: %s\n", businessUnitDetails.DisplayName)
			state.AppendToBuids(businessUnitDetails.DisplayName)
			// fmt.Printf("Business Unit ID: %s\n", bu.ID)
			// for _, link := range bu.Links {
			// 	fmt.Printf("Link: %s, Rel: %s, Method: %s\n", link.Href, link.Rel, link.Method)
			// }
		}

		// Here you would typically store the token securely
		// For now, just respond with success
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true}`))
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}

// /auth#access_token=tpa-f15f4140ae0a487b86be547fd76f&token_type=bearer&expires_in=360000
// CreateAuthPage returns an HTML page that extracts access token from URL fragment
// and sends it to the server before redirecting to Discord
func CreateAuthPage() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <title>Discord Authentication</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            text-align: center;
            margin: 50px;
            background-color: #36393f;
            color: #ffffff;
        }
        .container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #2f3136;
            border-radius: 8px;
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
        }
        h1 {
            color: #7289da;
        }
        #status {
            margin: 20px 0;
            padding: 10px;
            border-radius: 4px;
        }
        .success {
            background-color: #43b581;
        }
        .error {
            background-color: #f04747;
        }
        .loading {
            background-color: #faa61a;
        }
        button {
            background-color: #7289da;
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 16px;
            margin-top: 20px;
        }
        button:hover {
            background-color: #677bc4;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>Discord Authentication</h1>
        <div id="status" class="loading">Processing your authentication...</div>
        <div id="token-info"></div>
        <button id="redirect-btn" style="display:none;">Continue to Discord</button>
    </div>

    <script>
        document.addEventListener('DOMContentLoaded', function() {
            const statusEl = document.getElementById('status');
            const tokenInfoEl = document.getElementById('token-info');
            const redirectBtn = document.getElementById('redirect-btn');
            
            // Function to extract hash parameters
            function getHashParams() {
                const hash = window.location.hash.substring(1);
                const params = {};
                
                if (!hash) {
                    return params;
                }
                
                hash.split('&').forEach(pair => {
                    const [key, value] = pair.split('=');
                    params[key] = decodeURIComponent(value || '');
                });
                
                return params;
            }
            
            // Extract token from URL fragment
            const params = getHashParams();
            const accessToken = params['access_token'];
            const tokenType = params['token_type'];
            const expiresIn = params['expires_in'];
            
            if (!accessToken) {
                statusEl.textContent = 'Error: No access token found in URL';
                statusEl.className = 'error';
                return;
            }
            
            // Send token to server
            fetch('/store-token', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    access_token: accessToken,
                    token_type: tokenType,
                    expires_in: expiresIn
                })
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to store token');
                }
                return response.json();
            })
            .then(data => {
                statusEl.textContent = 'Authentication successful!';
                statusEl.className = 'success';
                tokenInfoEl.innerHTML = '<p>Your token has been securely stored.</p>';
                
                // Show redirect button
                redirectBtn.style.display = 'inline-block';
                redirectBtn.addEventListener('click', function() {
                    window.location.href = 'https://discord.com/app';
                });
                
                // Auto redirect after 5 seconds
                setTimeout(() => {
                    window.location.href = 'https://discord.com/app';
                }, 5000);
            })
            .catch(error => {
                console.error('Error:', error);
                statusEl.textContent = 'Error: ' + error.message;
                statusEl.className = 'error';
            });
        });
    </script>
</body>
</html>
`
}

// http get request
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

func processUserResponse(responseBody string) (*UserResponse, error) {
	var userResponse UserResponse
	err := json.Unmarshal([]byte(responseBody), &userResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing user response: %w", err)
	}

	return &userResponse, nil
}

func processBusinessUnitsResponse(responseBody string) (*BusinessUnitsResponse, error) {
	var businessUnitsResponse BusinessUnitsResponse
	err := json.Unmarshal([]byte(responseBody), &businessUnitsResponse)
	if err != nil {
		return nil, fmt.Errorf("error parsing business units response: %w", err)
	}

	return &businessUnitsResponse, nil
}

func processBusinessUnitDetails(responseBody string) (*BusinessUnitDetails, error) {
	var businessUnitDetails BusinessUnitDetails
	err := json.Unmarshal([]byte(responseBody), &businessUnitDetails)
	if err != nil {
		return nil, fmt.Errorf("error parsing business unit details: %w", err)
	}

	return &businessUnitDetails, nil
}
