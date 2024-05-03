package slotsuma

import (
	"Librarian/utils"

	"github.com/bwmarrin/discordgo"
)

func Slotsuma(s *discordgo.Session, guild string) ([]*discordgo.ApplicationCommand, error) {
	var commands []*discordgo.ApplicationCommand
	c, err := daily(s, guild)
	if err != nil {
		utils.Log.Error(err)
		return []*discordgo.ApplicationCommand{}, err
	}

	commands = append(commands, c...)

	c, err = coinflip(s, guild)
	if err != nil {
		utils.Log.Error(err)
		return []*discordgo.ApplicationCommand{}, err
	}

	return append(commands, c...), nil
}
