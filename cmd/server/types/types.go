package types

import "time"

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

type ReviewCreated struct {
	Events []struct {
		EventName string `json:"eventName"`
		Version   string `json:"version"`
		EventData struct {
			ID          string    `json:"id"`
			Language    string    `json:"language"`
			Title       string    `json:"title"`
			Text        string    `json:"text"`
			ReferenceID string    `json:"referenceId"`
			Stars       int       `json:"stars"`
			CreatedAt   time.Time `json:"createdAt"`
			IsVerified  bool      `json:"isVerified"`
			LocationID  string    `json:"locationId"`
			Link        string    `json:"link"`
			Consumer    struct {
				ID   string `json:"id"`
				Name string `json:"name"`
				Link string `json:"link"`
			} `json:"consumer"`
			Tags []struct {
				Group string `json:"group"`
				Value string `json:"value"`
			} `json:"tags"`
		} `json:"eventData"`
	} `json:"events"`
}
