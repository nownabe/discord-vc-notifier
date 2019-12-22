package main

import (
	"log"

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
