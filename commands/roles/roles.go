package roles

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"

	"Librarian/utils"
)

type roles struct {
	guild   string
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

var config map[string][]role_category = make(map[string][]role_category, 0)

var role_commands map[string]roles = make(map[string]roles, 0)

func init() {

	jsonFile, err := os.Open("./roles.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer jsonFile.Close()

	// Decode the JSON data into a struct
	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println(err)
		return
	}

	for k, v := range config {
		for _, category := range v {

			var options = make([]*discordgo.ApplicationCommandOption, 0)

			options = append(options, &discordgo.ApplicationCommandOption{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        category.Name,
				Description: category.Description,
			})

			role_commands[k] = roles{
				guild: k,
				command: &discordgo.ApplicationCommand{
					Name:        "roles",
					Description: "Sets your roles!",
					Options:     options,
				},
			}
		}
	}
}

func Roles(s *discordgo.Session, guild string) ([]*discordgo.ApplicationCommand, error) {

	var r = role_commands[guild]

	var err error
	r.command, err = s.ApplicationCommandCreate(s.State.User.ID, guild, r.command)

	if err != nil {
		return []*discordgo.ApplicationCommand{}, err
	}

	s.AddHandler(r.handler)
	s.AddHandler(r.role_select)

	return []*discordgo.ApplicationCommand{r.command}, nil

}

func (c *roles) handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {

		data := i.ApplicationCommandData()

		if data.ID == c.command.ID {

			var options = make([]discordgo.SelectMenuOption, 0)

			id, _ := get_category_id(c.guild, data.Options[0].Name)

			var category = config[c.guild][id]

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

		var categories = make([]string, 0, len(config[c.guild]))

		for _, category := range config[c.guild] {
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

			id, _ := get_category_id(c.guild, component.CustomID)

			//available roles
			var role_list []string
			for _, role := range config[c.guild][id].Roles {
				role_list = append(role_list, role.ID)
			}

			// Calculate roles to add and remove directly
			rolesToAdd := utils.Difference(component.Values, i.Member.Roles)
			rolesToRemove := utils.Difference(role_list, component.Values)

			// Assign new roles
			for _, role := range rolesToAdd {
				utils.Assign_Role(s, i.GuildID, i.Member.User.ID, role)
			}

			// Remove unselected roles
			for _, role := range rolesToRemove {
				utils.Remove_Role(s, i.GuildID, i.Member.User.ID, role)
			}
		}
	}
}

func get_category_id(guild string, name string) (int, error) {
	for i, category := range config[guild] {
		if name == category.Name {
			return i, nil
		}
	}
	return 0, fmt.Errorf("category \"%s\" does not exist", name)
}
