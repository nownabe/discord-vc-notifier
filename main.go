package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	cfg, err := newConfig()
	if err != nil {
		log.Fatalf("failed to get config: %v", err)
	}

	// Discord Bot

	discord, err := discordgo.New("Bot " + os.Getenv("DISCORD_BOT_TOKEN"))
	if err != nil {
		log.Fatalf("failed to build discord: %v", err)
	}

	notifiers := []notifier{}

	if cfg.DiscordChannelID != "" {
		notifiers = append(notifiers, &discordNotifier{
			channelID: cfg.DiscordChannelID,
			session:   discord,
		})
		log.Printf("added a discord notifier (%s)", cfg.DiscordChannelID)
	}

	discord.AddHandler(handler(notifiers))

	if err := discord.Open(); err != nil {
		log.Fatalf("failed to open discord: %v", err)
	}
	defer discord.Close()

	log.Printf("Opened discord")

	// HTTP Server

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "ok")
	})

	log.Printf("Listening on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, nil); err != nil {
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
