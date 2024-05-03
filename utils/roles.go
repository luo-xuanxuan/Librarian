package utils

import (
	"github.com/bwmarrin/discordgo"
)

func Assign_Role(s *discordgo.Session, guild string, user string, role string) error {
	err := s.GuildMemberRoleAdd(guild, user, role)
	if err != nil {
		Log.Error(err)
		return err
	}
	return nil
}

func Remove_Role(s *discordgo.Session, guild string, user string, role string) error {
	err := s.GuildMemberRoleRemove(guild, user, role)
	if err != nil {
		Log.Error(err)
		return err
	}
	return nil
}
