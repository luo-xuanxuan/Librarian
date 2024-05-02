package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
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

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}
	defer dg.Close()

	commands.Initialize_Handler(dg)
	defer commands.Unregister_Commands()

	create_commands()

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

func create_commands() {
	//misc commands
	//commands.Register_Command("", commands.Hi())
	//commands.Register_Command("", commands.Pipebomb())

	//roles
	// Open the JSON file
	jsonFile, err := os.Open("./roles.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	// Read the file content
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatalf("Error reading JSON file: %s", err)
	}

	// Use a map of string to json.RawMessage to dynamically handle the JSON structure
	var guilds map[string]json.RawMessage
	if err := json.Unmarshal(byteValue, &guilds); err != nil {
		log.Fatalf("Error decoding JSON: %s", err)
	}

	// Process each guild entry
	for guild, data := range guilds {
		fmt.Printf("Processing guild: %s\n", guild)
		commands.Register_Command(guild, commands.Roles(data))
	}

}
