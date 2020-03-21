package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("1 WORKER started")

	//dg, err := discordgo.New("Bot "+ os.Getenv("BOT_TOKEN"))
	dg, err := discordgo.New(os.Getenv("TOKEN"))
	if err != nil {
		fmt.Println("ERROR creating Discord session:", err)
		return
	}

	fmt.Println("2 Discord session created")

	dg.AddHandler(messageCreate)

	fmt.Println("3 Registred the messageCreate func")

	err = dg.Open()
	if err != nil {
		fmt.Println("ERROR opening connection:", err)
		return
	}

	fmt.Println("4 Opened a websocket")

	for i := 1; i > 0; i += 1 {
		time.Sleep(15 * time.Minute)
		fmt.Println("WORKER uped for", 15*i, "minutes")
	}
}

const (
	chJoraId = "635202206358044710"
	chBotId  = "467251523244523522"
	chBumpId = "569252448137510922"

	uStemaId = "522347439676588032"
	uSupId   = "569252448137510922"
	uBumpId  = "315926021457051650"

	gAhousId = "464116508252045312"
)

var (
	userId string
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == uStemaId &&
		m.ChannelID == chBotId &&
		m.Content == "R U TYT?" {
		s.ChannelMessageSend(chBotId, "E IM TYT!")
	}

	if m.Author.ID == uStemaId &&
		m.ChannelID == chBotId &&
		m.Content == "testSTEMAK" {
		fmt.Println("testSTEMAK runned")

		guild, _ := s.Guild(gAhousId)

		for _, member := range guild.Members {
			if member.User.String() == "stemak#2557" {
				s.ChannelMessageSend(chJoraId, ",add-money "+member.Mention()+" 10")
				break
			}
		}
	}

	if m.ChannelID == chBumpId &&
		m.Author.ID == uSupId &&
		len(m.Embeds) > 0 &&
		m.Embeds[0].Title == "Сервер Up" &&
		m.Embeds[0].Footer != nil {

		guild, _ := s.Guild(gAhousId)

		for _, member := range guild.Members {
			if member.User.String() == m.Embeds[0].Footer.Text {

				s.ChannelMessageSend(chJoraId, ",add-money "+member.Mention()+" 1000")
				s.ChannelMessageSend(chBumpId, member.Mention()+", Вы сделали Up сервера и Тихий Ужас вручил Вам 1000<:AH_AniCoin:579712087224483850>")
				fmt.Println("Sever uped by", m.Embeds[0].Footer.Text)
				break
			}
		}
	}
	
	if m.ChannelID == chBumpId &&
		m.Author.ID == uBumpId &&
		len(m.Embeds) > 0 {

		matched, _ := regexp.Match(`Server bumped by <@\d*>`, []byte(m.Embeds[0].Description))
		if matched {

			member := "<" + strings.Split(strings.Split(m.Embeds[0].Description, "<")[1], ">")[0] + ">"
			s.ChannelMessageSend(chJoraId, ",add-money "+member+" 1000")
			s.ChannelMessageSend(chBumpId, member+", Вы сделали Bump сервера и Скромный Модератор вручил Вам 1000<:AH_AniCoin:579712087224483850>")
			fmt.Println("Sever bumped by", m.Embeds[0].Footer.Text)
		}
	}
}
