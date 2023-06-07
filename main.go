package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"strings"
	"encoding/json"
	"io/ioutil"
	"time"
	"math/rand"
	"strconv"
	"database/sql"
	"regexp"

	"github.com/bwmarrin/discordgo"
	_ "github.com/lib/pq"
)

type Config struct {
	BotToken string `json:"token"`
	AppId string `json:"appId"`
	GuildId string `json:"guildId"`
	DbConnectionUrl string `json:"dbConnectionUrl"`
}

var db *sql.DB
var snowflakePattern *regexp.Regexp

func init() {
	rand.Seed(time.Now().UnixNano())

	snowflakePattern = regexp.MustCompile(`^\d{17,19}$`)
}

func main() {
	config, err := loadConfig()
	if err != nil {
		fmt.Println("error loading config:", err)
		return
	}

	db, err = sql.Open("postgres", config.DbConnectionUrl)
	if err != nil {
		panic(err)
	}

	dg, err := discordgo.New("Bot " + config.BotToken)
	if err != nil {
		fmt.Println("error creating Discord session:", err)
		return
	}

	err = registerCommands(dg, config.AppId, config.GuildId)
	if err != nil {
		fmt.Println("error setting slash commands:", err)
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(interactionCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening Discord connection:", err)
		return
	}

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
	db.Close()
}

func loadConfig() (*Config, error) {
	data, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

func registerCommands(s *discordgo.Session, appId string, guildId string) error {
	_, err := s.ApplicationCommandBulkOverwrite(appId, guildId, []*discordgo.ApplicationCommand {
		{
			Name: "ping",
			Description: "See if the bot is responding.",
		},
	})
	return err
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if strings.ToLower(m.Content) == "capy" {
		imgno := rand.Intn(747) + 1
		s.ChannelMessageSend(m.ChannelID, "https://api.capy.lol/v1/capybara/" + strconv.Itoa(imgno))
	
	} else if strings.HasPrefix(strings.ToLower(m.Content), "!time ") {
		userId := m.Content[6:]
		if !snowflakePattern.MatchString(userId) {
			s.ChannelMessageSend(m.ChannelID, "Invalid user ID.")
			return
		}
		var timezoneCell sql.NullString
		err := db.QueryRow("SELECT get_user_timezone($1);", userId).Scan(&timezoneCell)
		if err != nil {
			fmt.Println("error scanning db row:", err)
			return
		}
		if timezoneCell.Valid {
			s.ChannelMessageSend(m.ChannelID, "That user's timezone is **" + timezoneCell.String + "**.")
		} else {
			s.ChannelMessageSend(m.ChannelID, "There is no timezone set for that user.")
		}
	}
}

func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	data := i.ApplicationCommandData()
	switch data.Name {
	case "ping":
		err := s.InteractionRespond(
			i.Interaction,
			&discordgo.InteractionResponse {
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData {
					Content: "Pong!",
				},
			},
		)
		if err != nil {
			fmt.Println("error responding to interaction:", err)
		}
	}
}
