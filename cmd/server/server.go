package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/liukaku/discord-tp/cmd/server/handlers"
	"github.com/liukaku/discord-tp/cmd/server/handlers/utils"
	"github.com/liukaku/discord-tp/cmd/server/pages"
	"github.com/liukaku/discord-tp/cmd/server/types"
)

func CreateHttpServer(discord *discordgo.Session, state types.SharedState) {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("received get request")
		fmt.Println("basepath URL frag:", r.URL.Fragment)
		if r.Method != http.MethodGet {
			fmt.Println("Method received: ", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

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
	})

	http.HandleFunc("/trustpilot", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("received trustpilot request")
		if r.Method != http.MethodPost {
			fmt.Println("Method received: ", r.Method)
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		// parse out the request body and log it
		requestBody := r.Body
		bodyBytes, err := io.ReadAll(requestBody)
		if err != nil {
			fmt.Println("Error reading request body:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		fmt.Println("Request body:", string(bodyBytes))

		var trustpilotRequest types.ReviewCreated
		err = json.Unmarshal(bodyBytes, &trustpilotRequest)
		if err != nil {
			fmt.Println("Error parsing request body:", err)
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		fmt.Println("Parsed Trustpilot request:", len(trustpilotRequest.Events))
		ids := state.GetChannelIDs()
		for n := range ids {
			starString := ""
			for range trustpilotRequest.Events[0].EventData.Stars {
				starString += "⭐"
			}

			link := utils.ConvertTrustpilotApiUrlToPublic(trustpilotRequest.Events[0].EventData.Link)

			_, err := discord.ChannelMessageSendComplex(ids[n], &discordgo.MessageSend{
				Content: fmt.Sprintf("New Trustpilot review received:\n**%s**\n%s\nRating: %s\nLink: %s",
					trustpilotRequest.Events[0].EventData.Consumer.Name,
					trustpilotRequest.Events[0].EventData.Text,
					starString,
					link,
				),
				Embeds: []*discordgo.MessageEmbed{
					{
						Title:       "Trustpilot Review",
						Description: trustpilotRequest.Events[0].EventData.Text,
						URL:         link,
						Color:       0x00ff00, // Green color
						Fields: []*discordgo.MessageEmbedField{
							{
								Name:   "Author",
								Value:  trustpilotRequest.Events[0].EventData.Consumer.Name,
								Inline: true,
							},
						},
						// Thumbnail: &discordgo.MessageEmbedThumbnail{
						// 	URL: trustpilotRequest.Events[0].EventData.Consumer.Image,
						// },
					},
				},
			})
			if err != nil {
				fmt.Println("Error sending message:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			discord.ChannelMessageSend(ids[n], fmt.Sprintf("New Trustpilot review received: %s", trustpilotRequest.Events[0].EventData.Text))

		}
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
		htmlPage := pages.CreateAuthPage()
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

		handlers.StoreTokenHandler(w, r, state)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true}`))
	})

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
