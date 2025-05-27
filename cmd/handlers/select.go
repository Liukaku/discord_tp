package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var selectHandlersMap = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"select-1": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("You selected option 1 from %s", i.MessageComponentData().CustomID),
			},
		})
	},
	"select-2": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("You selected option 2 from %s", i.MessageComponentData().CustomID),
			},
		})
	},
}

func SelectHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if handler, ok := selectHandlersMap[i.MessageComponentData().CustomID]; ok {
		handler(s, i)
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unknown select option",
			},
		})
	}
}
