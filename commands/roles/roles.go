package roles

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"

	"Librarian/config"
	"Librarian/utils"
)

type roles struct {
	command *discordgo.ApplicationCommand
}

type role_category struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	MinValues   int    `json:"MinValues"`
	MaxValues   int    `json:"MaxValues"`
	Roles       []role `json:"Roles"`
}

type role struct {
	Name string `json:"Name"`
	ID   string `json:"ID"`
}

var cfg map[string][]role_category = make(map[string][]role_category, 0)

var role_commands map[string]roles = make(map[string]roles, 0)

func init() {
	config.Load_Config("roles", "roles", &cfg)
}

func build_guild(s *discordgo.Session, guild string) ([]*discordgo.ApplicationCommand, error) {

	var com = roles{
		command: &discordgo.ApplicationCommand{
			Name:        "roles",
			Description: "Sets your roles!",
		},
	}

	var err error
	com.command, err = s.ApplicationCommandCreate(s.State.User.ID, guild, com.command)

	if err != nil {
		return []*discordgo.ApplicationCommand{}, err
	}

	role_commands[guild] = com

	s.AddHandler(com.handler)
	s.AddHandler(com.category_select)
	s.AddHandler(com.role_select)

	return []*discordgo.ApplicationCommand{com.command}, nil
}

func Roles(s *discordgo.Session, guild string) ([]*discordgo.ApplicationCommand, error) {

	r, err := build_guild(s, guild)
	if err != nil {
		return []*discordgo.ApplicationCommand{}, err
	}

	/*
		c, err := Role_Edit(s, guild)
		if err != nil {
			return []*discordgo.ApplicationCommand{}, err
		}*/

	return r, nil

}

func (c *roles) handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()

	if data.ID != c.command.ID {
		return
	}

	var options = make([]discordgo.SelectMenuOption, 0)

	for k, v := range cfg[i.GuildID] {
		options = append(options,
			discordgo.SelectMenuOption{
				Label:       v.Name,
				Value:       fmt.Sprintf("%d", k),
				Description: v.Description,
			})
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    fmt.Sprintf("roles%s", i.GuildID),
					Placeholder: "Select a category",
					Options:     options,
				},
			},
		},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    "Please choose from the options:",
			Components: components,
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		utils.Log.Error(err)
	}
}

func (c *roles) category_select(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	// Type assertion to check it's a ComponentInteraction
	component := i.MessageComponentData()

	if component.CustomID != fmt.Sprintf("roles%s", i.GuildID) {
		return
	}

	index, err := strconv.ParseInt(component.Values[0], 10, 64)
	if err != nil {
		utils.Log.Error(err)
		return
	}

	category := cfg[i.GuildID][index]

	var options = make([]discordgo.SelectMenuOption, 0)

	for _, role := range category.Roles {
		options = append(options, discordgo.SelectMenuOption{
			Label:   role.Name,
			Value:   role.ID,
			Default: utils.Contains(i.Member.Roles, role.ID),
		})
	}

	max_value := category.MaxValues
	if max_value > len(options) {
		max_value = len(options)
	}

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					CustomID:    category.Name,
					Placeholder: "Select your role(s)",
					MinValues:   &category.MinValues,
					MaxValues:   max_value,
					Options:     options,
				},
			},
		},
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    category.Description,
			Components: components,
			Flags:      discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		utils.Log.Error(err)
	}

}

func (c *roles) role_select(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	// Type assertion to check it's a ComponentInteraction
	component := i.MessageComponentData()

	var categories []string

	for _, category := range cfg[i.GuildID] {
		categories = append(categories, category.Name)
	}

	if !utils.Contains(categories, component.CustomID) {
		return
	}

	// Acknowledge the interaction without sending a message
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})
	if err != nil {
		fmt.Println("Failed to respond to interaction:", err)
		return
	}

	id, _ := get_category_id(i.GuildID, component.CustomID)

	//available roles
	var role_list []string
	for _, role := range cfg[i.GuildID][id].Roles {
		role_list = append(role_list, role.ID)
	}

	// Calculate roles to add and remove directly
	roles_to_add := utils.Difference(component.Values, i.Member.Roles)
	roles_to_remove := utils.Difference(role_list, component.Values)
	roles_to_remove = utils.Intersection(roles_to_remove, i.Member.Roles)

	guild_roles, err := s.GuildRoles(i.GuildID)
	if err != nil {
		utils.Log.Error(err)
	}

	roles_map := make(map[string]string)
	for _, gr := range guild_roles {
		roles_map[gr.ID] = gr.Name
	}

	// Assign new roles
	for j, role := range roles_to_add {
		utils.Assign_Role(s, i.GuildID, i.Member.User.ID, role)
		roles_to_add[j] = roles_map[role]
	}
	utils.Log.WithField("User", i.Member.User.Username).WithField("Guild", i.GuildID).WithField("Roles", strings.Join(roles_to_add, ", ")).Info("Added Roles")

	// Remove unselected roles
	for j, role := range roles_to_remove {
		utils.Remove_Role(s, i.GuildID, i.Member.User.ID, role)
		roles_to_remove[j] = roles_map[role]
	}
	utils.Log.WithField("User", i.Member.User.Username).WithField("Guild", i.GuildID).WithField("Roles", strings.Join(roles_to_remove, ", ")).Info("Removed Roles")
}

func get_category_id(guild string, name string) (int, error) {
	for i, category := range cfg[guild] {
		if name == category.Name {
			return i, nil
		}
	}
	return 0, fmt.Errorf("category \"%s\" does not exist", name)
}
