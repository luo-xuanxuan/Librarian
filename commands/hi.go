package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type hi_command struct {
	command discordgo.ApplicationCommand
}

func Hi() *hi_command {

	var cmd = discordgo.ApplicationCommand{
		Name:        "hello",
		Description: "Says hello!",
	}

	return &hi_command{
		command: cmd,
	}
}

func (c hi_command) get_command() *discordgo.ApplicationCommand {
	return &c.command
}

func (c hi_command) handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
