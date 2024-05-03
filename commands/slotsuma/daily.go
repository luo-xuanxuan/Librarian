package slotsuma

import (
	"Librarian/user"
	"Librarian/utils"
	"fmt"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/bwmarrin/discordgo"
)

var daily_command *discordgo.ApplicationCommand

var daily_amount = 1000
var streak_amount = 100

var daily_handled bool = false

func init() {
	daily_command = &discordgo.ApplicationCommand{
		Name:        "daily",
		Description: "Collect your daily gil!!",
	}
}

func daily(s *discordgo.Session, guild string) ([]*discordgo.ApplicationCommand, error) {
	var err error
	daily_command, err = s.ApplicationCommandCreate(s.State.User.ID, guild, daily_command)
	if err != nil {
		utils.Log.Error(err)
		return []*discordgo.ApplicationCommand{}, err
	}

	if !daily_handled {
		s.AddHandler(daily_handler)
		daily_handled = true
	}

	return []*discordgo.ApplicationCommand{daily_command}, nil
}

func daily_handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
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

	last_daily := active_user.Get("slotsuma_daily_timer")
	if last_daily == nil {
		last_daily = 0
	}

	now := time.Now().Unix()

	time_passed := now - last_daily.(int64)

	if time_passed < 86400 {
		// a day hasnt passed yet

		response := fmt.Sprintf("Sorry you already claimed your daily recently!!! T^T\nYou can claim again in <t:%d:R>", (last_daily.(int64) + 86400))

		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: response,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		if err != nil {
			utils.Log.Error(err)
		}
		return
	}

	streak := active_user.Get("slotsuma_daily_streak")
	if (streak == nil) || (time_passed > 172800) {
		streak = 0
	}

	gil := active_user.Get("gil")
	if gil == nil {
		gil = 0
	}

	gil = gil.(int64) + 1000 + (streak.(int64) * int64(streak_amount))

	active_user.Set("gil", gil.(int64))
	active_user.Set("slotsuma_daily_streak", (streak.(int64) + 1))
	active_user.Set("slotsuma_daily_timer", now)
	p := message.NewPrinter(language.English)

	response := p.Sprintf("You collected your daily %dg!", daily_amount)
	if streak.(int64) > 0 {
		response = p.Sprintf("%s\nBonus +%dg for %d streak!", response, streak.(int64)*int64(streak_amount), streak.(int64))
	}
	response = p.Sprintf("%s\nYou total gil right now is: %dg", response, gil.(int64))

	err := s.InteractionRespond(
		i.Interaction,
		&discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
	if err != nil {
		utils.Log.Error(err)
	}

}
