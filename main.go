package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"Librarian/commands"

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
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	//dg.AddHandler(messageCreate)

	//dg.AddHandler(pipebomb)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	commands.Initialize_Handler(dg)

	create_commands()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	commands.Unregister_Commands()

	// Cleanly close down the Discord session.
	dg.Close()
}

func create_commands() {
	commands.Register_Command("", commands.Hi())
	commands.Register_Command("", commands.Pipebomb())
}
