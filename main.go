package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

const prefix string = "!8b"
const version string = "1.1.0"

func goDotEnvVariable(key string) string {
	// Load .env file.
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Return value from key provided.
	return os.Getenv(key)
}

func main() {
	// Grab bot token env var.
	botToken := goDotEnvVariable("BOT_TOKEN")

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + botToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// guildID := m.Message.GuildID

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Grab message content from guild.
	content := m.Content

	if strings.Contains(content, prefix+"help") {
		// Grab author
		author := m.Author.Username

		// commandHelpTitle := "Looks like you need a hand. Check out my goodies below... \n \n"
		greeting := "Hi " + author + "! \n \n"
		introduction := "8Ball bot is here to help.\n"

		// // Notes
		note1 := "This bot will provides a carefully determined answer to your desired question 🎱\n \n"
		// note2 := "❕ Commands are case-sensitive. Lower-case only :) \n \n"
		note3 := "\n👨🏼‍💻 Dev: Narsiq#5638. DM me for requests/questions/sups\n \n"

		commandPrefix := prefix

		// // Commands
		commandsHeader := "\nCOMMANDS:\n\n"
		commandHelpMessage := "❔  " + commandPrefix + "help - Provides a list of my commands. \n\n"
		commandAnswer := "🎱  " + commandPrefix + " - Returns a carefully determined answer to your desired question.\n\n"
		commandInvite := "🔗  " + commandPrefix + "invite - A invite link for the 8Ball Bot. \n\n"
		commandSite := "🔗  " + commandPrefix + "site - Link to the 8Ball website. \n\n"
		commandSupport := "💝  " + commandPrefix + "support - Link to the 8Ball Patreon. \n\n"
		commandStats := "📊  " + commandPrefix + "stats - Check out 8Ball stats. \n\n"
		commandVersion := "🤖  " + commandPrefix + "version - Current 8Ball version. \n\n"

		messageFull := greeting + introduction + note1 + commandsHeader + commandAnswer + commandHelpMessage + commandInvite + commandSite + commandSupport + commandStats + commandVersion + note3

		message := "```\n" + messageFull + "\n```"

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, prefix+"site") {
		// Build start vote message
		author := m.Author.Username
		message := "Here ya go " + author + "..." + "\n\n" + "https://discordbots.dev/"

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, prefix+"support") {
		// Build start vote message
		author := m.Author.Username
		message := "Thanks for thinking of me " + author + " 💖." + "\n\n" + "https://www.patreon.com/BotVoteTo"

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, prefix+"version") {
		// Build start vote message
		message := "8Ball is currently running version " + version

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.Contains(content, prefix+"stats") {
		// TODO: This will need to be updated to iterate through
		// all shards once the bot joins 1,000 servers.
		// var guilds []string
		guilds := s.State.Ready.Guilds
		fmt.Println(len(guilds))
		guildCount := len(guilds)
		guildCountStr := strconv.Itoa(guildCount)

		fmt.Println(guilds)
		fmt.Printf("t1: %T\n", guilds)

		// // Build start vote message
		message := "8Ball bot is currently on " + guildCountStr + " servers! Noice..."

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.EqualFold(content, prefix+"invite") {
		author := m.Author.Username

		// // Build start vote message
		message := "Wow! Such nice " + author + ". Thanks for spreading the 💖. Here is an invite link made just for you... \n \n" + "https://discord.com/api/oauth2/authorize?client_id=979170113842479205&permissions=100416&scope=bot"

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, message, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}

	if strings.EqualFold(content, prefix) {
		// Call answer func to generate values.
		answer := getAnswer()

		// Grab author
		author := m.Author.Username

		// Build start vote message
		messageGreet := author + "... \n"
		messageAnswer := "```fix" + "\n" + "🎱 " + answer + "\n" + "```"
		messageFull := messageGreet + messageAnswer

		// Send start vote message
		_, err := s.ChannelMessageSendReply(m.ChannelID, messageFull, m.Reference())
		if err != nil {
			fmt.Println(err)
		}
	}
}

func getAnswer() string {
	csvFile, err := os.Open("answers.csv")
	if err != nil {
		fmt.Println(err)
	}
	defer csvFile.Close()

	// Read csv file
	csvLines, err := csv.NewReader(csvFile).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	// Generate random number using min/max index of csv file lines.
	rand.Seed(time.Now().UnixNano())
	min := 0
	max := 20
	randomIndex := rand.Intn(max-min+1) + min
	result := csvLines[randomIndex]

	return result[0]
}
