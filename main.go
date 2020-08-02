package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	stopMonitor := monitorMemStats()
	defer stopMonitor()

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

	if cfg.SlackWebhookURL != "" {
		sn := newSlackNotifier(cfg.SlackWebhookURL, cfg.SlackChannel,
			cfg.SlackUsername, cfg.SlackIconEmoji)
		notifiers = append(notifiers, sn)
		log.Printf("added a slack notifier")
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

		var name string
		if m.Nick != "" {
			name = m.Nick
		} else {
			name = m.User.Username
		}

		msg := fmt.Sprintf("%sがボイチャにいるよ", name)

		for _, n := range ns {
			n.notify(msg)
		}
	}
}

func monitorMemStats() func() {
	s := new(runtime.MemStats)
	t := time.NewTicker(time.Minute)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				return
			case <-t.C:
				runtime.ReadMemStats(s)
				json, err := json.Marshal(s)
				if err == nil {
					log.Printf("%s", json)
				} else {
					log.Printf("Failed to marshal mem stats: %v", err)
				}
			}
		}
	}()

	log.Printf("started monitoring mem stats")

	return func() {
		t.Stop()
		done <- struct{}{}
	}
}
