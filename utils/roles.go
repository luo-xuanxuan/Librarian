package utils

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func Assign_Role(s *discordgo.Session, guild string, user string, role string) error {
	err := s.GuildMemberRoleAdd(guild, user, role)
	if err != nil {
		fmt.Printf("Failed to assign role: %v\n", err)
		return err
	}
	fmt.Println("Role assigned successfully")
	return nil
}

func Remove_Role(s *discordgo.Session, guild string, user string, role string) error {
	err := s.GuildMemberRoleRemove(guild, user, role)
	if err != nil {
		fmt.Printf("Failed to remove role: %v\n", err)
		return err
	}
	fmt.Println("Role removed successfully")
	return nil
}
