package main

import (
	"github.com/bwmarrin/discordgo"
)

var (
	registeredCommands []*discordgo.ApplicationCommand = make([]*discordgo.ApplicationCommand, 0)
)

func registerCommand(s *discordgo.Session, guild string, cmd *discordgo.ApplicationCommand) {
	c, err := s.ApplicationCommandCreate(s.State.User.ID, guild, cmd)

	if err != nil {
		println("Cannot create '%v' command: %v", cmd.Name, err)
	}
	registeredCommands = append(registeredCommands, c)
}

func unregisterCommands(s *discordgo.Session) {
	for _, v := range registeredCommands {
		err := s.ApplicationCommandDelete(s.State.User.ID, v.GuildID, v.ID)
		if err != nil {
			println("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}
