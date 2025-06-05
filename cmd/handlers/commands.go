package handlers

import (
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func createDropdowns(i *discordgo.InteractionCreate, state SharedState) *discordgo.InteractionResponseData {
	buids := state.GetBuids()
	var dropdowns []discordgo.SelectMenuOption
	fmt.Println("Creating dropdowns for business units:", buids)
	fmt.Println("Creating dropdowns for business units:", len(buids))
	if len(buids) == 0 {
		dropdowns = []discordgo.SelectMenuOption{
			{
				Label:       "No Business Units",
				Value:       "no-buids",
				Description: "No business units available",
				Default:     true,
			},
		}
	} else {
		dropdowns = make([]discordgo.SelectMenuOption, len(buids))
		for i, buid := range buids {
			dropdowns[i] = discordgo.SelectMenuOption{
				Label:       buid,
				Value:       buid,
				Description: "Business Unit: " + buid,
				Default:     false,
			}
		}
	}

	return &discordgo.InteractionResponseData{
		Content: "Lets take a look at these settings" + i.Interaction.Member.User.ID,
		Flags:   discordgo.MessageFlagsEphemeral,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "bu-select",
						Placeholder: "Select a business unit",
						Options:     dropdowns,
						MenuType:    discordgo.StringSelectMenu,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "channel-select",
						Placeholder: "Select a channel",
						MenuType:    discordgo.ChannelSelectMenu,
						ChannelTypes: []discordgo.ChannelType{
							discordgo.ChannelTypeGuildText,
						},
					},
				},
			},
		},
	}
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, state SharedState){
	"settings": func(s *discordgo.Session, i *discordgo.InteractionCreate, state SharedState) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: createDropdowns(i, state),
		})

		if err != nil {
			fmt.Println("Error responding to interaction: ", err)
			return
		}
	},
	"test": func(s *discordgo.Session, i *discordgo.InteractionCreate, state SharedState) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Test command executed",
			},
		})
	},
	"modal": func(s *discordgo.Session, i *discordgo.InteractionCreate, state SharedState) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Test command executed",
			},
		})
	},
	"login": func(s *discordgo.Session, i *discordgo.InteractionCreate, state SharedState) {
		// respond with a button to open a website
		fmt.Println("Login command executed by:", i.Interaction.Member.User.Username)
		baseUrl := "https://authenticate.tp-staging.com/?redirect_uri="
		redirectUrl := "http://localhost:8080/auth" // Replace with your redirect URL
		// load in client id from env
		clientId := os.Getenv("CLIENT_ID")
		queryParams := fmt.Sprintf("&client_id=%s&response_type=token", clientId)
		serverIdQueryParams := fmt.Sprintf("&guild_id=%s", i.Interaction.GuildID)
		fullUrl := baseUrl + redirectUrl + queryParams + serverIdQueryParams

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Login to the bot",
				Flags:   discordgo.MessageFlagsEphemeral,
				Components: []discordgo.MessageComponent{
					discordgo.ActionsRow{
						Components: []discordgo.MessageComponent{
							discordgo.Button{
								Label: "Login",
								Style: discordgo.LinkButton,
								URL:   fullUrl, // Replace with your login URL
							},
						},
					},
				},
			},
		})
		if err != nil {
			fmt.Println("Error responding to login command:", err)
			return
		}
	},
}

func CommandHandler(s *discordgo.Session, i *discordgo.InteractionCreate, state SharedState) {
	fmt.Printf("command: '%s' from %s\n", i.ApplicationCommandData().Name, i.Interaction.Member.User.Username)
	if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
		h(s, i, state)
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unknown Command",
			},
		})
	}
}
