package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	dg.AddHandler(pipebomb)

	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Type == discordgo.InteractionApplicationCommand {
			if i.ApplicationCommandData().Name == "hello" {
				err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
					Type: discordgo.InteractionResponseChannelMessageWithSource,
					Data: &discordgo.InteractionResponseData{
						Content: "Hello there!",
						Components: []discordgo.MessageComponent{
							discordgo.ActionsRow{
								Components: []discordgo.MessageComponent{
									discordgo.Button{
										Emoji: &discordgo.ComponentEmoji{
											Name: "📜",
										},
										Label: "Documentation",
										Style: discordgo.LinkButton,
										URL:   "https://discord.com/developers/docs/interactions/message-components#buttons",
									},
									discordgo.Button{
										Emoji: &discordgo.ComponentEmoji{
											Name: "🔧",
										},
										Label: "Discord developers",
										Style: discordgo.LinkButton,
										URL:   "https://discord.gg/discord-developers",
									},
									discordgo.Button{
										Emoji: &discordgo.ComponentEmoji{
											Name: "🦫",
										},
										Label: "Discord Gophers",
										Style: discordgo.LinkButton,
										URL:   "https://discord.gg/7RuRrVHyXF",
									},
								},
							},
						},
						Flags: discordgo.MessageFlagsEphemeral,
					},
				})
				if err != nil {
					fmt.Printf("Cannot send message: %v\n", err)
				}
				message, err := s.InteractionResponse(i.Interaction)
				if err != nil {
					fmt.Printf("Cannot fetch interaction message: %v\n", err)
					return
				}

				// Now you can access the message ID
				fmt.Printf("Message ID: %s\n", message.ID)

				newMessageContent := "bye"

				_, err = s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
					Content: &newMessageContent,
				})
				if err != nil {
					fmt.Printf("Failed to edit interaction response: %v\n", err)
					return
				}
			}
		}
	})

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	createCommands(dg)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	unregisterCommands(dg)

	// Cleanly close down the Discord session.
	dg.Close()
}

func createCommands(s *discordgo.Session) {
	registerCommand(s, "491232530314297365", &discordgo.ApplicationCommand{
		Name:        "hello",
		Description: "Says hello!",
	})
	registerCommand(s, "", &pipebombCommand)
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	//fmt.Println(m.Message.Content)
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
		s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{Content: "Pong!", Reference: m.Reference(), Flags: discordgo.MessageFlagsEphemeral})
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
