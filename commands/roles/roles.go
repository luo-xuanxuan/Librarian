package roles

import (
	"github.com/bwmarrin/discordgo"

	"Librarian/config"
)

type Role_Command struct {
	command *discordgo.ApplicationCommand
}

type role_category struct {
	Name        string            `json:"Name"`
	Description string            `json:"Description"`
	MinValues   int               `json:"MinValues"`
	MaxValues   int               `json:"MaxValues"`
	Roles       map[string]string `json:"Roles"` //Role Name to Role ID
}

var cfg map[string][]role_category = make(map[string][]role_category, 0)

var handled bool = false

func init() {
	//Load roles from json
	//TODO: Rework config loading to make more sense
	config.Load_Config("roles", "roles", &cfg)
}

func (r *Role_Command) Preload(s *discordgo.Session, guild string) error {
	if !handled {
		s.AddHandler(category_select_handler)
		s.AddHandler(role_select_handler)
	}
	return nil
}

func (r *Role_Command) Reference() *discordgo.ApplicationCommand {
	if r.command == nil {
		r.command = &discordgo.ApplicationCommand{
			Name:        "roles",
			Description: "Sets your roles!",
		}
	}
	return r.command
}

func (r *Role_Command) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) error {

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				create_category_select(i.GuildID),
			},
		},
	}

	return s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Please choose from the options:",
			Components: components,
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
}
