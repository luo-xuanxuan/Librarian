package roles

import (
	"Librarian/config"
	"Librarian/utils"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var edit_commands map[string]*discordgo.ApplicationCommand = make(map[string]*discordgo.ApplicationCommand, 0)

func Role_Edit(s *discordgo.Session, guild string) ([]*discordgo.ApplicationCommand, error) {

	var command = &discordgo.ApplicationCommand{
		Name:        "role-edit",
		Description: "Edit optional roles lists",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "create",
				Description: "Create new role selection.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "label",
						Description: "Label of role selection.",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "description",
						Description: "Description of role selection.",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "min_value",
						Description: "Minimum selection.",
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "max_value",
						Description: "Maximum selection.",
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "delete",
				Description: "Deletes a role selection.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "label",
						Description: "Label of role selection to delete.",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "add",
				Description: "Adds role to selection.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "label",
						Description: "Label of role selection to append.",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "Role to add to selection.",
						Required:    true,
					},
				},
			},
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "remove",
				Description: "Removes role from selection.",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:        discordgo.ApplicationCommandOptionString,
						Name:        "label",
						Description: "Label of role selection to remove role from.",
						Required:    true,
					},
					{
						Type:        discordgo.ApplicationCommandOptionRole,
						Name:        "role",
						Description: "Role to remove from selection.",
						Required:    true,
					},
				},
			},
		},
	}

	var err error
	command, err = s.ApplicationCommandCreate(s.State.User.ID, guild, command)

	if err != nil {
		return []*discordgo.ApplicationCommand{}, err
	}

	edit_commands[guild] = command

	s.AddHandler(role_edit)

	return []*discordgo.ApplicationCommand{command}, nil
}

func role_edit(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//early return if interaction is incorrect type
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()

	//early return if the ID does not match
	if data.ID != edit_commands[i.GuildID].ID {
		return
	}

	//early return if user lacks required permission. ilegal :(
	if i.Member.Permissions&discordgo.PermissionManageRoles != discordgo.PermissionManageRoles {
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "You lack ManageRoles Permission :(",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		if err != nil {
			utils.Log.Error(err)
		}
		return
	}

	//utils.Log.Info(data.Options[0].Name)

	switch data.Options[0].Name {
	case "create":
		create(s, i)
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("%s selection created!", data.Options[0].Options[0].StringValue()),
					Flags:   discordgo.MessageFlagsEphemeral, //Ephemeral flag keeps message invisible to other users
				},
			})
		if err != nil {
			utils.Log.Error(err)
		}
	case "delete":
		delete_roles(s, i)
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("%s selection removed!", data.Options[0].Options[0].StringValue()),
					Flags:   discordgo.MessageFlagsEphemeral, //Ephemeral flag keeps message invisible to other users
				},
			})
		if err != nil {
			utils.Log.Error(err)
		}
	case "add":
		add(s, i)
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("%s added to %s!", data.Options[0].Options[1].RoleValue(s, i.GuildID).Name, data.Options[0].Options[0].StringValue()),
					Flags:   discordgo.MessageFlagsEphemeral, //Ephemeral flag keeps message invisible to other users
				},
			})
		if err != nil {
			utils.Log.Error(err)
		}
	case "remove":
		remove(s, i)
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("%s removed from %s!", data.Options[0].Options[1].StringValue(), data.Options[0].Options[0].StringValue()),
					Flags:   discordgo.MessageFlagsEphemeral, //Ephemeral flag keeps message invisible to other users
				},
			})
		if err != nil {
			utils.Log.Error(err)
		}
	default:
		utils.Log.Error("Invalid Role Edit Subcommand")
	}

	config.Save_Config("roles", "roles", cfg)
}

func create(s *discordgo.Session, i *discordgo.InteractionCreate) {

	var data = i.ApplicationCommandData()

	var min_value = 0
	var max_value = 0
	for _, op := range data.Options[0].Options {
		if op.Name == "min_value" {
			min_value = int(op.IntValue())
		}
		if op.Name == "max_value" {
			max_value = int(op.IntValue())
		}
	}

	cfg[i.GuildID] = append(cfg[i.GuildID],
		role_category{
			Name:        data.Options[0].Options[0].StringValue(),
			Description: data.Options[0].Options[1].StringValue(),
			MinValues:   min_value,
			MaxValues:   max_value,
		},
	)
}

func delete_roles(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	for j, v := range cfg[i.GuildID] {
		if v.Name == data.Options[0].Options[0].StringValue() {
			// Remove the item from the slice
			cfg[i.GuildID] = append(cfg[i.GuildID][:j], cfg[i.GuildID][j+1:]...)
			// Since only one item can have the name, we return here if this is within a function
			return
		}
	}

}

func add(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	label := data.Options[0].Options[0].StringValue()
	r := data.Options[0].Options[1].RoleValue(s, i.GuildID)

	for j, _ := range cfg[i.GuildID] {
		if cfg[i.GuildID][j].Name == label {
			cfg[i.GuildID][j].Roles = append(cfg[i.GuildID][j].Roles, role{
				Name: r.Name,
				ID:   r.ID,
			})
			return
		}
	}
}

func remove(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()

	label := data.Options[0].Options[0].StringValue()
	r := data.Options[0].Options[1].RoleValue(s, i.GuildID)

	for j, v := range cfg[i.GuildID] {
		if v.Name == label {
			v.Roles = append(v.Roles, role{
				Name: r.Name,
				ID:   r.ID,
			})
			v.Roles = append(v.Roles[:j], v.Roles[j+1:]...)
			return
		}
	}
}
