package main

import (
	"fmt"
	"time"
	"os"

	"github.com/bwmarrin/discordgo"
)

func main() {
	fmt.Println("1 WORKER started")
	
	dg, err := discordgo.New(os.Getenv("EMAIL"), os.Getenv("PASS"))
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
		fmt.Println("WORKER uped for", 15 * i, "minutes")
	}
}

var (
	isntProcess bool = true
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == "522347439676588032" &&
		m.ChannelID == "467251523244523522" && 
		m.Content == "R U TYT?" {
		  s.ChannelMessageSend("467251523244523522", "E IM TYT!")
	}

	if m.Author.ID == "522347439676588032" &&
		m.ChannelID == "467251523244523522" && 
		m.Content == "test1" {
		  s.ChannelMessageSend("635202206358044710", ",add-money \"ДобраяKnopKa | 雨 🌧#2575\" 13")
	}

	if m.Author.ID == "522347439676588032" &&
		m.ChannelID == "467251523244523522" && 
		m.Content == "test2" {
			s.ChannelMessageSend("635202206358044710", ",add-money \"GLuK | AniHouseTV#0015\" 13")
	}

	if m.Author.ID == "522347439676588032" &&
		m.ChannelID == "467251523244523522" && 
		m.Content == "test3" {
			s.ChannelMessageSend("635202206358044710", ",add-money \"Турианка#1704\" 13")
	}

	if len(m.Embeds) > 0 && 
		m.ChannelID == "569252448137510922" && 
		m.Author.ID == "464272403766444044" && 
		m.Embeds[0].Title == "Сервер Up" && 
		m.Embeds[0].Footer != nil {

		s.ChannelMessageSend("635202206358044710", ",add-money \"" + m.Embeds[0].Footer.Text + "\" 1000")
		s.ChannelMessageSend("569252448137510922", m.Embeds[0].Footer.Text + ", за Up сервера Вам вам были начислены 1000 <:AH_AniCoin:579712087224483850>")

		fmt.Println("Sever uped by", m.Embeds[0].Footer.Text)
	}
}
