package commands

import (
	"encoding/json"
	"fmt"

	"github.com/bwmarrin/discordgo"

	"Librarian/utils"
)

type roles struct {
	command    discordgo.ApplicationCommand
	categories []role_category
}

type role struct {
	Name string `json:"Name"`
	ID   string `json:"ID"`
}

type role_category struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	MinValues   int    `json:"MinValues"`
	MaxValues   int    `json:"MaxValues"`
	Roles       []role `json:"Roles"`
}

func Roles(data json.RawMessage) *roles {

	// Decode the JSON data into the roles struct

	var categories []role_category

	var _ = json.Unmarshal(data, &categories)

	var options = make([]*discordgo.ApplicationCommandOption, 0)

	for _, category := range categories {
		options = append(options, &discordgo.ApplicationCommandOption{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        category.Name,
			Description: category.Description,
		})
	}

	var r = roles{
		command: discordgo.ApplicationCommand{
			Name:        "roles",
			Description: "Sets your roles!",
			Options:     options,
		},
		categories: categories,
	}

	session.AddHandler(r.role_select)

	return &r

}

func (c *roles) get_command() *discordgo.ApplicationCommand {
	return &c.command
}

func (c *roles) set_command(cmd *discordgo.ApplicationCommand) {
	c.command = *cmd
}

func (c *roles) get_category_id(name string) int {
	for i, category := range c.categories {
		if name == category.Name {
			return i
		}
	}
	return -1
}

func (c *roles) handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {

		data := i.ApplicationCommandData()

		if data.ID == c.command.ID {

			var options = make([]discordgo.SelectMenuOption, 0)

			var id = c.get_category_id(data.Options[0].Name)

			var category = c.categories[id]

			for _, role := range category.Roles {
				options = append(options, discordgo.SelectMenuOption{
					Label:   role.Name,
					Value:   role.ID,
					Default: utils.Contains(i.Member.Roles, role.ID),
				})
			}

			components := []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.SelectMenu{
							CustomID:    category.Name,
							Placeholder: "Select your role(s)",
							MinValues:   &category.MinValues,
							MaxValues:   category.MaxValues,
							Options:     options,
						},
					},
				},
			}

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content:    "Please choose from the options:",
					Components: components,
					Flags:      discordgo.MessageFlagsEphemeral,
				},
			})
		}
	}
}

func (c *roles) role_select(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionMessageComponent {
		// Type assertion to check it's a ComponentInteraction
		component := i.MessageComponentData()

		var categories = make([]string, 0, len(c.categories))

		for _, category := range c.categories {
			categories = append(categories, category.Name)
		}

		if utils.Contains(categories, component.CustomID) {

			// Acknowledge the interaction without sending a message
			err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseDeferredMessageUpdate,
			})
			if err != nil {
				fmt.Println("Failed to respond to interaction:", err)
				return
			}

			var id = c.get_category_id(component.CustomID)

			//available roles
			var role_list []string
			for _, role := range c.categories[id].Roles {
				role_list = append(role_list, role.ID)
			}

			// Calculate roles to add and remove directly
			rolesToAdd := utils.Difference(component.Values, i.Member.Roles)
			rolesToRemove := utils.Difference(role_list, component.Values)

			// Assign new roles
			for _, role := range rolesToAdd {
				assign_role(s, i.GuildID, i.Member.User.ID, role)
			}

			// Remove unselected roles
			for _, role := range rolesToRemove {
				remove_role(s, i.GuildID, i.Member.User.ID, role)
			}
		}
	}
}

func assign_role(s *discordgo.Session, guild string, user string, role string) error {
	err := s.GuildMemberRoleAdd(guild, user, role)
	if err != nil {
		fmt.Printf("Failed to assign role: %v\n", err)
		return err
	}
	fmt.Println("Role assigned successfully")
	return nil
}

func remove_role(s *discordgo.Session, guild string, user string, role string) error {
	err := s.GuildMemberRoleRemove(guild, user, role)
	if err != nil {
		fmt.Printf("Failed to remove role: %v\n", err)
		return err
	}
	fmt.Println("Role removed successfully")
	return nil
}
