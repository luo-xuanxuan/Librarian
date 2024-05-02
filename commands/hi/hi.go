package hi

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var name string = "hello"
var handled bool = false

func Hi(s *discordgo.Session, guild string) ([]*discordgo.ApplicationCommand, error) {

	//Define our command details
	command := &discordgo.ApplicationCommand{
		Name:        name,
		Description: "Says hello!",
	}

	//Register the command to the guild
	var err error
	command, err = s.ApplicationCommandCreate(s.State.User.ID, guild, command)

	if err != nil {
		return []*discordgo.ApplicationCommand{}, err
	}

	//we only want one handler instance, but we might need multiple commands registered, so we just check if we handled.
	if !handled {
		s.AddHandler(hi)
		handled = true
	}

	//return the command reference so we can remove it on shutdown
	return []*discordgo.ApplicationCommand{command}, nil
}

// hi command handler
func hi(s *discordgo.Session, i *discordgo.InteractionCreate) {

	//early return if interaction is incorrect type
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	//early return if the name does not match
	if i.ApplicationCommandData().Name != name {
		return
	}

	//respond to interaction with a message
	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Hey there!",
				Flags:   discordgo.MessageFlagsEphemeral, //Ephemeral flag keeps message invisible to other users
			},
		})
	if err != nil {
		fmt.Printf("Cannot send message: %v\n", err)
	}
}
