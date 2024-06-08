package command_handler

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Reference() *discordgo.ApplicationCommand
	Handler(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

type PreloadCommand interface {
	Preload(s *discordgo.Session, guild string) error
}

var Stop_Handler func()

var command_references map[string]Command

func init() {
	command_references = make(map[string]Command)
}

func Register_Command(s *discordgo.Session, guild string, commands ...Command) {
	for _, cmd := range commands {
		if preload, ok := cmd.(PreloadCommand); ok {
			preload.Preload(s, guild)
		}
		app := cmd.Reference()
		app, err := s.ApplicationCommandCreate(s.State.User.ID, guild, app)
		if err != nil {
			log.Println(err.Error())
			return
		}
		command_references[app.ID] = cmd
	}
}

func Start_Handler(s *discordgo.Session) {
	Stop_Handler = s.AddHandler(handler)
}

func handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	err := command_references[i.ApplicationCommandData().ID].Handler(s, i)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func Unregister_Commands(s *discordgo.Session) {
	Stop_Handler()
	for _, v := range command_references {
		s.ApplicationCommandDelete(s.State.User.ID, v.Reference().GuildID, v.Reference().ID)
	}
}
