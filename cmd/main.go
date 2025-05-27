package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/liukaku/discord-tp/cmd/handlers"
)

var discord *discordgo.Session

var (
	GuildID        = flag.String("guild", "", "Test guild ID. If not passed - bot registers commands globally")
	BotToken       = flag.String("token", "", "Bot access token")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "settings",
		Description: "Open the settings modal",
		Type:        discordgo.ChatApplicationCommand,
	},
	{
		Name:        "test",
		Description: "Test command",
		Type:        discordgo.ChatApplicationCommand,
	},
	{
		Name:        "modal",
		Description: "Test command",
		Type:        discordgo.ChatApplicationCommand,
	},
}

func bodyOne(i *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		Content: "Lets take a look at these settings" + i.Interaction.Member.User.ID,
		Flags:   discordgo.MessageFlagsEphemeral,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.SelectMenu{
						CustomID:    "select-1",
						Placeholder: "Select an option",
						Options: []discordgo.SelectMenuOption{
							{
								Label:       "Option 1",
								Value:       "option-1",
								Description: "This is option 1",
								Default:     false,
							},
							{
								Label:       "Option 2",
								Value:       "option-2",
								Description: "This is option 2",
								Default:     false,
							},
						},
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

func bodyTwo(i *discordgo.InteractionCreate) *discordgo.InteractionResponseData {
	return &discordgo.InteractionResponseData{
		CustomID: "modals_survey_" + i.Interaction.Member.User.ID,
		Title:    "Modals survey",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:    "opinion",
						Label:       "What is your opinion on them?",
						Style:       discordgo.TextInputShort,
						Placeholder: "Don't be shy, share your opinion with us",
						Required:    true,
						MaxLength:   300,
						MinLength:   10,
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.TextInput{
						CustomID:  "suggestions",
						Label:     "What would you suggest to improve them?",
						Style:     discordgo.TextInputParagraph,
						Required:  false,
						MaxLength: 2000,
					},
				},
			},
		},
	}
}

var commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
	"settings": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: bodyOne(i),
		})

		if err != nil {
			fmt.Println("Error responding to interaction: ", err)
			return
		}
	},
	// "settings": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// 	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
	// 		Type: discordgo.InteractionResponseModal,
	// 		Data: bodyTwo(i),
	// 	})

	// 	if err != nil {
	// 		fmt.Println("Error responding to interaction: ", err)
	// 		return
	// 	}
	// },
	"test": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Test command executed",
			},
		})
	},
	"modal": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Test command executed",
			},
		})
	},
	"ready": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Test command executed",
			},
		})
	},
}

func addHandlers() {
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			fmt.Println("New Event: InteractionCreate: ")
			fmt.Println(i.ApplicationCommandData().Name)
			fmt.Println(i.ApplicationCommandData().Options)
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			fmt.Println("New Event: InteractionCreate: ")
			fmt.Println(i.MessageComponentData().CustomID)
			fmt.Println(i.MessageComponentData().Values)
			handlers.SelectHandler(s, i)
		}
	})
}

func main() {

	go createHttpServer()

	fmt.Println("Bot Launching")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}

	discordKey := os.Getenv("DISCORD_TOKEN")
	fmt.Println(discordKey)

	discord, err = discordgo.New("Bot " + discordKey)
	if err != nil {
		fmt.Println("Error creating Discord session, make sure you have the correct token")
		panic(err)
	}

	err = discord.Open()
	// Register the command handlers
	addHandlers()
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))
	for i, v := range commands {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, *GuildID, v)
		if err != nil {
			fmt.Printf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer discord.Close()

	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages
	// This function will be called (due to AddHandler above) when the bot receives
	// the "ready" event from Discord.
	discord.AddHandler(ready)

	// This function will be called (due to AddHandler above) every time a new
	// message is created on any channel that the autenticated bot has access to.
	discord.AddHandler(messageCreate)

	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	discord.Close()

	if *RemoveCommands {
		fmt.Printf("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

		for _, v := range registeredCommands {
			err := discord.ApplicationCommandDelete(discord.State.User.ID, *GuildID, v.ID)
			if err != nil {
				fmt.Println("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	fmt.Printf("Gracefully shutting down.")
}

func test(m *discordgo.MessageCreate) {
	channel, err := discord.Channel("1234")
	if err != nil {
		fmt.Println("Error getting channel: ", err)
		return
	}
	discord.MessageThreadStartComplex(channel.ID, m.ID, &discordgo.ThreadStart{
		Name: "Test thread",
	})
}

func ready(s *discordgo.Session, event *discordgo.Ready) {

	fmt.Println("New Event: Ready: ")
	// parse the event to get the guilds
	// fmt.Print(json.Marshal(event))
	fmt.Println(event.Guilds)
	for _, guild := range event.Guilds {
		fmt.Printf("Guild ID: %s, Guild Name: %s\n", *GuildID, guild.Name)
		for _, channel := range guild.Channels {
			fmt.Printf("Channel ID: %s, Channel Name: %s\n", channel.ID, channel.Name)
			discord.ChannelMessage(channel.ID, "Hello from the bot!")
		}
	}

}

func modalCreate(s *discordgo.Session, m *discordgo.InteractionCreate) {
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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println("New Event: Message: ")
	fmt.Println(m.Message.Content)
	fmt.Println(m.Message.ChannelID)
	//create a modal for the user to fill out
	s.ChannelMessageSendReply(m.ChannelID, "bing bong", m.Message.Reference())
}

type HttpResponse struct {
	Hi  int `json:"hi"`
	Bye int `json:"bye"`
}

func createHttpServer() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		panic(err)
	}
	// channelId := os.Getenv("CHANNEL_ID")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// wg.Add(1)
		fmt.Println("received get request")
		fmt.Println(r.Method)
		fmt.Println(r.Body)
		// read request body
		var resp map[string]interface{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "can't read body", http.StatusBadRequest)
			panic(err)
		}
		json.Unmarshal(body, &resp)
		fmt.Println(string(body))
		fmt.Println(resp)
		// io.WriteString(w, fmt.Sprintf("Received request: %s \n song number: %s", r.Method, k))
		// get all guilds the bot is in
		guilds, err := discord.UserGuilds(100, "", "", true)
		for i, guild := range guilds {
			fmt.Printf("Guild %d: %s (%s)\n", i, guild.Name, guild.ID)
			channels, _ := discord.GuildChannels(guild.ID)
			for i, channel := range channels {
				fmt.Printf("Channel %d: %s (%s)\n", i, channel.Name, channel.ID)
				discord.ChannelMessageSend(channel.ID, "iyaaaaaaa")
			}
		}
		// get all channels in the guild

	})
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
