package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"Librarian/command_handler"
	"Librarian/commands/roles"
	"Librarian/commands/universalis"

	"github.com/bwmarrin/discordgo"
)

// Variables used for command line parameters
var (
	Token string
)

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	// Create a new Discord session using the provided bot token.
	s, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	s.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer s.Close()

	command_handler.Start_Handler(s)
	defer command_handler.Unregister_Commands(s)

	create_commands(s)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func create_commands(s *discordgo.Session) {
	//roles
	command_handler.Register_Command(s, "781419076462837760", &roles.Role_Command{}, &universalis.Universalis_Command{})
}
