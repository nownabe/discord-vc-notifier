package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

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

type slackNotifier struct {
	channel    string
	client     *http.Client
	iconEmoji  string
	username   string
	webhookURL string
}

type slackMessage struct {
	Channel   string `json:"channel,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"`
	Text      string `json:"text"`
	Username  string `json:"username,omitempty"`
}

func newSlackNotifier(webhookURL, channel, username, iconEmoji string) notifier {
	return &slackNotifier{
		channel:    channel,
		client:     &http.Client{},
		iconEmoji:  iconEmoji,
		username:   username,
		webhookURL: webhookURL,
	}
}

func (n *slackNotifier) notify(msg string) {
	m := slackMessage{
		Channel:   n.channel,
		IconEmoji: n.iconEmoji,
		Text:      msg,
		Username:  n.username,
	}

	json, err := json.Marshal(m)
	if err != nil {
		log.Printf("failed to marshal slack message: %v", err)
		return
	}

	req, err := http.NewRequest("POST", n.webhookURL, bytes.NewReader(json))
	if err != nil {
		log.Printf("failed to create new request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := n.client.Do(req)
	if err != nil {
		log.Printf("failed to send message to slack: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read response body: %v", err)
		return
	}

	if resp.StatusCode >= 400 {
		log.Printf("slack webhook request failed: %s (%d)", body, resp.StatusCode)
		return
	}
}
