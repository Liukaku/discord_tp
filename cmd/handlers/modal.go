package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func ModalCreate(s *discordgo.Session, m *discordgo.InteractionCreate) {
	fmt.Println("New Event: Message: ")
	fmt.Println(m.Interaction.ChannelID)

	s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			Title:    "Modal Title",
			CustomID: "modal-1",
			Components: []discordgo.MessageComponent{
				discordgo.StringSelectMenu: &discordgo.SelectMenu{
					CustomID:    "select-1",
					Placeholder: "Select an option",
					Options: []discordgo.SelectMenuOption{
						{
							Label:       "Option 1",
							Value:       "option-1",
							Description: "This is option 1",
						},
						{
							Label:       "Option 2",
							Value:       "option-2",
							Description: "This is option 2",
						},
					},
				},
			},
		},
	})
}
