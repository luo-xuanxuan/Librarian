package roles

import (
	"Librarian/utils"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func create_category_select(guild string) *discordgo.SelectMenu {

	var category_options = make([]discordgo.SelectMenuOption, len(cfg[guild]))

	for i, v := range cfg[guild] {
		category_options[i] = discordgo.SelectMenuOption{
			Label:       v.Name,
			Value:       v.Name,
			Description: v.Description,
		}
	}

	return &discordgo.SelectMenu{
		CustomID:    "role_category_select",
		Placeholder: "Select a Role Category.",
		Options:     category_options,
	}

}

func category_select_handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionMessageComponent {
		return
	}

	component := i.MessageComponentData()

	if component.CustomID != "role_category_select" {
		return
	}

	var category *role_category
	var index int
	for idx, v := range cfg[i.GuildID] {
		utils.Log.WithField("value", v.Name).Info()
		if component.Values[0] == v.Name {
			category = &v
			index = idx
			break
		}
	}

	role_select := create_role_select(category, i)

	role_select.CustomID = fmt.Sprintf("%s_%d", role_select.CustomID, index)

	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				role_select,
			},
		},
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
