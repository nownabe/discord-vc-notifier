package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

type notifier interface {
	notify(string)
}

type discordNotifier struct {
	channelID string
	session   *discordgo.Session
}

func (n *discordNotifier) notify(msg string) {
	if _, err := n.session.ChannelMessageSend(n.channelID, msg); err != nil {
		log.Printf("failed to send message to discord: %v", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("failed to build discord: %v", err)
	}

	notifiers := []notifier{
		&discordNotifier{channelID: "658184345382551563", session: discord},
	}

	discord.AddHandler(handler(notifiers))

	if err := discord.Open(); err != nil {
		log.Fatalf("failed to open discord: %v", err)
	}
	defer discord.Close()

	log.Printf("Opened discord")
	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
	log.Printf("Terminating...")
}

func handler(ns []notifier) func(*discordgo.Session, *discordgo.VoiceStateUpdate) {
	return func(s *discordgo.Session, vsu *discordgo.VoiceStateUpdate) {
		if vsu.VoiceState.ChannelID == "" {
			return
		}

		m, err := s.GuildMember(vsu.VoiceState.GuildID, vsu.VoiceState.UserID)
		if err != nil {
			log.Printf("failed to get member: %v", err)
			return
		}

		msg := fmt.Sprintf("%sがボイチャにいるよ", m.Nick)

		for _, n := range ns {
			n.notify(msg)
		}
	}
}
