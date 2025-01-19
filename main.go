package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var discord *discordgo.Session

func main() {

	go createHttpServer()

	fmt.Println("Bot Launching")
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	discordKey := os.Getenv("DISCORD_TOKEN")
	fmt.Println(discordKey)

	discord, err = discordgo.New("Bot " + discordKey)
	if err != nil {
		panic(err)
	}

	discord.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMessages
	// This function will be called (due to AddHandler above) when the bot receives
	// the "ready" event from Discord.
	discord.AddHandler(ready)

	// This function will be called (due to AddHandler above) every time a new
	// message is created on any channel that the autenticated bot has access to.
	discord.AddHandler(messageCreate)

	err = discord.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
	}
	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)
	<-sigch

	discord.Close()
}

func ready(s *discordgo.Session, event *discordgo.Ready) {

	fmt.Println("New Event: Ready: ")
	fmt.Println(event)

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	fmt.Println("New Event: Message: ")
	fmt.Println(m.Message.Content)
	fmt.Println(m.Message.ChannelID)
	s.ChannelMessageSendReply(m.ChannelID, "bing bong", m.Message.Reference())
}

type HttpResponse struct {
	Hi  int `json:"hi"`
	Bye int `json:"bye"`
}

func createHttpServer() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	channelId := os.Getenv("CHANNEL_ID")

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
		discord.ChannelMessageSend(channelId, string(body))

	})
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
