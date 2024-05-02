package commands

import (
	"Librarian/utils"
	"fmt"
	"image"
	"log"
	"net"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/bwmarrin/discordgo"
)

var server_address string = "192.168.1.136:4210"
var name = "pipebomb"
var handled bool = false

func Pipebomb(s *discordgo.Session, guild string) ([]*discordgo.ApplicationCommand, error) {

	command := &discordgo.ApplicationCommand{
		Name:        name,
		Description: "Puts a message on the pipebomb!",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "text",
				Description: "The message for the screen.",
				Required:    true,
			},
		},
	}

	//Register the command to the guild
	var err error
	command, err = s.ApplicationCommandCreate(s.State.User.ID, guild, command)

	if err != nil {
		return []*discordgo.ApplicationCommand{}, err
	}

	//we only want one handler instance, but we might need multiple commands registered, so we just check if we handled.
	if !handled {
		s.AddHandler(pipebomb)
		handled = true
	}

	//return the command reference so we can remove it on shutdown
	return []*discordgo.ApplicationCommand{command}, nil
}

func pipebomb(s *discordgo.Session, i *discordgo.InteractionCreate) {

	//early return if interaction is incorrect type
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()

	//early return if the name does not match
	if data.Name != name {
		return
	}

	userInput := ""
	for _, option := range data.Options {
		if option.Name == "text" {
			userInput = option.StringValue()
		}
	}

	userInput = utils.Sanitize_Discord_Text(userInput)

	response := fmt.Sprintf("You said: \"%s\"", userInput)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		fmt.Printf("Cannot send message: %v\n", err)
	}

	println(i.Member.User.Username + ": " + userInput)

	buf := text_to_bytes(i.Member.User.Username + ": " + userInput)

	// Resolve UDP address
	addr, err := net.ResolveUDPAddr("udp", server_address)
	if err != nil {
		log.Fatal(err)
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write(buf)
	if err != nil {
		log.Fatal(err)
	}
}

func text_to_bytes(text string) []byte {
	maxWidth := 128  // Image width
	lineHeight := 13 // Line height in pixels based on the chosen font
	maxLines := 64 / lineHeight

	// Create a new blank image 128x64
	img := image.NewGray(image.Rect(0, 0, maxWidth, 64))

	// Drawer with a basic font
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.White,
		Face: basicfont.Face7x13,
	}

	// Word wrap logic
	words := strings.Fields(text)
	var lines []string
	line := ""

	for _, word := range words {
		if drawer.MeasureString(line+" "+word).Ceil() > maxWidth {
			if line != "" {
				lines = append(lines, line)
				line = ""
			}

			wordWidth := drawer.MeasureString(word).Ceil()
			if wordWidth > maxWidth {
				// Split the word into parts that fit within maxWidth
				for len(word) > 0 {
					cut := len(word)
					for drawer.MeasureString(word[:cut]).Ceil() > maxWidth && cut > 1 {
						cut--
					}
					lines = append(lines, word[:cut])
					word = word[cut:]
				}
			} else {
				line = word // Start a new line with the word
			}
		} else {
			if line == "" {
				line = word
			} else {
				line += " " + word
			}
		}
	}
	if line != "" {
		lines = append(lines, line)
	}

	// Draw the lines
	for i, line := range lines {
		if i >= maxLines {
			break
		}
		drawer.Dot = fixed.P(0, (i+1)*lineHeight) // Positioning text correctly
		drawer.DrawString(line)
	}

	// Prepare a buffer to store bits, 128x64 bits = 1024 bytes
	buf := make([]byte, 1024)

	// Convert image pixels to bits
	for y := 0; y < 64; y++ {
		for x := 0; x < 128; x++ {
			index := y*128 + x
			if img.GrayAt(x, y).Y > 128 { // Assuming a threshold for "on" pixel
				buf[index/8] |= 1 << (7 - uint(index)%8)
			}
		}
	}

	return buf
}
