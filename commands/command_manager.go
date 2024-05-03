package commands

import (
	"fmt"

	"github.com/bwmarrin/discordgo"

	"Librarian/commands/hi"
	"Librarian/commands/roles"
)

var (
	registered_commands []*discordgo.ApplicationCommand                                                      = make([]*discordgo.ApplicationCommand, 0)
	available_packages  map[string]func(*discordgo.Session, string) ([]*discordgo.ApplicationCommand, error) = make(map[string]func(*discordgo.Session, string) ([]*discordgo.ApplicationCommand, error), 0)
)

func init() {
	available_packages["Roles"] = roles.Roles
	available_packages["Hi"] = hi.Hi
}

func Clean(s *discordgo.Session) {
	unregister_commands(s)
}

func Register_Packages(s *discordgo.Session, guild string, packages ...string) {
	for _, pkg := range packages {
		commands, err := available_packages[pkg](s, guild)
		if err != nil {
			fmt.Println(err)
			continue
		}
		register_commands(commands...)
	}
}

func register_commands(commands ...*discordgo.ApplicationCommand) {
	registered_commands = append(registered_commands, commands...)
}

func unregister_commands(s *discordgo.Session) {
	for _, v := range registered_commands {
		err := s.ApplicationCommandDelete(s.State.User.ID, v.GuildID, v.ID)
		if err != nil {
			println("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}
