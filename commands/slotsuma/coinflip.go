package slotsuma

import (
	"Librarian/user"
	"Librarian/utils"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
)

var coinflip_command *discordgo.ApplicationCommand

var coinflip_handled bool = false

func init() {

	rand.Seed(time.Now().UnixNano())

	coinflip_command = &discordgo.ApplicationCommand{
		Name:        "coinflip",
		Description: "Bet on a 50/50 to double your bet!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "bet",
				Description: "Amount to bet.",
			},
		},
	}
}

func coinflip(s *discordgo.Session, guild string) ([]*discordgo.ApplicationCommand, error) {
	var err error
	coinflip_command, err = s.ApplicationCommandCreate(s.State.User.ID, guild, coinflip_command)
	if err != nil {
		utils.Log.Error(err)
		return []*discordgo.ApplicationCommand{}, err
	}

	if !coinflip_handled {
		s.AddHandler(coinflip_handler)
		coinflip_handled = true
	}

	return []*discordgo.ApplicationCommand{coinflip_command}, nil
}

func coinflip_handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//early return if interaction is incorrect type
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	//early return if the name does not match
	if i.ApplicationCommandData().Name != daily_command.Name {
		return
	}

	var active_user = user.User{}
	active_user.Load(i.User.ID)

	gil := active_user.Get("gil")
	if gil == nil {
		gil = 0
	}

	var bet int64 = 0

	if len(i.ApplicationCommandData().Options) > 0 {
		bet = i.ApplicationCommandData().Options[0].IntValue()
		if bet < 0 {
			//dont bet negative bitch
			return
		}
		if gil.(int64) < bet {
			//poor bitch
			return
		}
	}

	var random int64 = int64(rand.Intn(2))

	//var embed *discordgo.MessageEmbed

	if random == 1 {
		gil = gil.(int64) + bet
	}

	if random == 0 {
		gil = gil.(int64) - bet
	}

}
