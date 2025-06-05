package handlers

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// Define interface for the shared state
type SharedState interface {
	AppendToStateArr(values ...string)
	AppendToChannelIDs(values ...string)
	GetBuids() []string
}

var selectHandlersMap = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, state SharedState){
	"select-1": func(s *discordgo.Session, i *discordgo.InteractionCreate, state SharedState) {
		state.AppendToStateArr(i.MessageComponentData().Values...)
		fmt.Println("Selected values:", i.MessageComponentData().Values)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("You selected option 1 from %s", i.MessageComponentData().Values),
			},
		})
	},
	"channel-select": func(s *discordgo.Session, i *discordgo.InteractionCreate, state SharedState) {
		state.AppendToChannelIDs(i.MessageComponentData().Values...)
		fmt.Println("Selected values:", i.MessageComponentData().Values)

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   discordgo.MessageFlagsEphemeral,
				Content: fmt.Sprintf("You selected option 2 from %s", i.MessageComponentData().CustomID),
			},
		})
	},
}

func SelectHandler(s *discordgo.Session, i *discordgo.InteractionCreate, state SharedState) {
	if handler, ok := selectHandlersMap[i.MessageComponentData().CustomID]; ok {
		handler(s, i, state)
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Unknown select option",
			},
		})
	}
}
