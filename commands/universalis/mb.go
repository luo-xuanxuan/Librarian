package universalis

import "github.com/bwmarrin/discordgo"

type Universalis_Command struct {
	command *discordgo.ApplicationCommand
}

func (uc *Universalis_Command) Reference() *discordgo.ApplicationCommand {
	if uc.command == nil {
		uc.command = &discordgo.ApplicationCommand{
			Name:        "mb",
			Description: "Check Universalis for item data!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "region-dc-world",
					Description: "Select an option",
					Required:    true,
					Choices: []*discordgo.ApplicationCommandOptionChoice{
						{
							Name:  "North America",
							Value: "na",
						},
						{
							Name:  "Europe",
							Value: "eu",
						},
						{
							Name:  "Oceania",
							Value: "oce",
						},
					},
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "fuzzy item name",
					Required:    true,
				},
			},
		}
	}
	return uc.command
}

func (uc *Universalis_Command) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	return nil
}
