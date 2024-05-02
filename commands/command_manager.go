package commands

import (
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	get_command() *discordgo.ApplicationCommand
	handler(s *discordgo.Session, i *discordgo.InteractionCreate)
}

var (
	registered_commands []*discordgo.ApplicationCommand = make([]*discordgo.ApplicationCommand, 0)
	session             *discordgo.Session
)

func Initialize_Handler(s *discordgo.Session) {
	session = s
}

func Register_Command(guild string, command Command) error {

	c, err := session.ApplicationCommandCreate(session.State.User.ID, guild, command.get_command())

	if err != nil {
		return err
	}
	session.AddHandler(command.handler)

	registered_commands = append(registered_commands, c)

	return nil
}

func Unregister_Commands() {
	for _, v := range registered_commands {
		err := session.ApplicationCommandDelete(session.State.User.ID, v.GuildID, v.ID)
		if err != nil {
			println("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}
