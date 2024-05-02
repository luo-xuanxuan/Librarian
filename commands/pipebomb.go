package commands

import (
	"fmt"
	"image"
	"log"
	"net"
	"regexp"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"

	"github.com/bwmarrin/discordgo"
)

type pipebomb struct {
	command discordgo.ApplicationCommand
}

var server_address = "192.168.1.136:4210"

func Pipebomb() *pipebomb {
	return &pipebomb{
		command: discordgo.ApplicationCommand{
			Name:        "Pipebomb",
			Description: "Puts a message on the pipebomb!",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "text",
					Description: "The message for the screen.",
					Required:    true,
				},
			},
		},
	}
}

func (c pipebomb) get_command() *discordgo.ApplicationCommand {
	return &c.command
}

func (c pipebomb) handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type == discordgo.InteractionApplicationCommand {

		data := i.ApplicationCommandData()

		if data.Name == "pipebomb" {

			userInput := ""
			for _, option := range data.Options {
				if option.Name == "text" {
					userInput = option.StringValue()
				}
			}

			userInput = sanitizeDiscordText(userInput)

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

			buf := textToBytes(i.Member.User.Username + ": " + userInput)

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
	}
}

func sanitizeDiscordText(input string) string {
	// Regular expression to find custom Discord emojis which are in format <:name:id>
	re := regexp.MustCompile(`<:([^:]+):\d+>`)
	// Replace all instances of the custom emoji with just its name wrapped in colons
	return re.ReplaceAllString(input, ":$1:")
}

func textToBytes(text string) []byte {
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
