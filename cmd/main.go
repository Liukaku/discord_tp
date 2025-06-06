package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	"github.com/liukaku/discord-tp/cmd/handlers"
	"github.com/liukaku/discord-tp/cmd/server"
)

var discord *discordgo.Session

// Create a struct to hold our shared state with a mutex
type SharedState struct {
	sync.RWMutex
	StateArr      []string
	ChannelIDs    []string
	BusinessUnits []string
}

// Create a global instance of our shared state
var sharedState = SharedState{
	StateArr:      []string{},
	ChannelIDs:    []string{},
	BusinessUnits: []string{},
}

// Add helper methods to safely access and modify state
func (s *SharedState) GetStateArr() []string {
	s.RLock()
	defer s.RUnlock()
	// Return a copy to prevent race conditions
	result := make([]string, len(s.StateArr))
	copy(result, s.StateArr)
	return result
}

func (s *SharedState) GetChannelIDs() []string {
	s.RLock()
	defer s.RUnlock()
	result := make([]string, len(s.ChannelIDs))
	copy(result, s.ChannelIDs)
	return result
}

func (s *SharedState) GetBuids() []string {
	s.RLock()
	defer s.RUnlock()
	result := make([]string, len(s.BusinessUnits))
	copy(result, s.BusinessUnits)
	return result
}

func (s *SharedState) AppendToStateArr(values ...string) {
	s.Lock()
	defer s.Unlock()
	s.StateArr = append(s.StateArr, values...)
	fmt.Println("Updated State Array:", s.StateArr)
}

func (s *SharedState) AppendToChannelIDs(values ...string) {
	s.Lock()
	defer s.Unlock()
	s.ChannelIDs = append(s.ChannelIDs, values...)
	fmt.Println("Updated Channel IDs:", s.ChannelIDs)
}

func (s *SharedState) AppendToBuids(values ...string) {
	s.Lock()
	defer s.Unlock()
	s.BusinessUnits = append(s.BusinessUnits, values...)
	fmt.Println("Updated Business Units:", s.BusinessUnits)
}

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
	{
		Name:        "login",
		Description: "Login to the bot",
		Type:        discordgo.ChatApplicationCommand,
	},
}

func addHandlers() {
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			fmt.Println("New Command Event: InteractionCreate: ")
			fmt.Println(i.ApplicationCommandData().Name)
			fmt.Println(i.ApplicationCommandData().Options)
			handlers.CommandHandler(s, i, &sharedState)
		case discordgo.InteractionMessageComponent:
			fmt.Println("New Interaction Event: InteractionCreate: ")
			fmt.Println(i.MessageComponentData().CustomID)
			fmt.Println(i.MessageComponentData().Values)
			handlers.SelectHandler(s, i, &sharedState)
		}
	})
}

func main() {

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

	// Create a new HTTP server to handle requests
	go server.CreateHttpServer(discord, &sharedState)

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

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println("New Event: Message: ")
	fmt.Println(m.Message.Content)
	fmt.Println(m.Message.ChannelID)
	fmt.Println(m.ReferencedMessage.Embeds[0].URL)
	//create a modal for the user to fill out
	s.ChannelMessageSendReply(m.ChannelID, "bing bong", m.Message.Reference())
}

type HttpResponse struct {
	Hi  int `json:"hi"`
	Bye int `json:"bye"`
}
