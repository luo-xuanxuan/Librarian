package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type hi struct {
	command discordgo.ApplicationCommand
}

func Hi() *hi {

	return &hi{
		command: discordgo.ApplicationCommand{
			Name:        "hello",
			Description: "Says hello!",
		},
	}
}

func (c hi) get_command() *discordgo.ApplicationCommand {
	return &c.command
}

func (c hi) handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {

		data := i.ApplicationCommandData()

		if data.Name == c.command.Name {
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Hey there!",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			if err != nil {
				fmt.Printf("Cannot send message: %v\n", err)
			}
		}
	}
}
