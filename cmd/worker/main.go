package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	fmt.Println("1 WORKER started")

	dg, err := discordgo.New(os.Getenv("TOKEN"))
	errorLog(err, "ERROR creating Discord session:")

	fmt.Println("2 Discord session created")

	dg.AddHandler(messageCreate)

	fmt.Println("3 Registred the messageCreate func")

	err = dg.Open()
	errorLog(err, "ERROR opening connection:")

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
	uSupId   = "464272403766444044"
	uBumpId  = "315926021457051650"

	gAhousId = "464116508252045312"
)

var (
	userId string

	responces = []string{
		"Вы сделали %s сервера и Тихий Ужас вручил Вам %s",
		"Вы героически поймали %s и Вас наградили %s",
		"за %s сервера Жмяк отдала Вам свои печеньки и вы получили %s",
		"после кибератаки вы подняли сервер своим %sом и получили %s",
		"пожертвовав жизнью на войне за %s, Вас посмертно наградили %s",
		"смакуючи ельфійською абракадаброю (%s), ви начарували %s",
		"сыр съел сырный сырник %sая сервер, Вам заплатили моральную компенсацию в %s",
		"Вы помогли Ведьмаку с %sом, за что Вам заплатили __ЧЕКАННЫМИ__ %s",
		"Вы сделали %s сервера и Скромный Модератор вручил Вам %s",
		"Вы съели свою свою поджелудочную во время %sа и нашли в ней %s",
		"Вы собрали с тысячи людей по АниКоину и %sнули сервер. Вы получили %s",
		"Вы попытались выписать бан Кнопычу, но сделали %s и получили %s",
		"на вечерних посиделках с Хикаро вы сражались за %s. Хикаро Вас наградил %s",
		"Вы сохранили последние пять минут и угрожающе сделали %s. Глюк расстрогался и отдал Вам %s",
		"Вы наблюдали за программированием Стёмы и Меро, но не забыли сделать %s и получили %s",
	}
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == uStemaId &&
		m.ChannelID == chJoraId &&
		m.Content == "R U TYT?" {
		s.ChannelMessageSend(chJoraId, "E IM TYT!")
	}

	if m.Author.ID == uStemaId &&
		m.ChannelID == chJoraId &&
		m.Content == "testSTEMAK" {

		guild, err := s.Guild(gAhousId)
		errorLog(err, "ERROR get guild failure:")

		for _, member := range guild.Members {
			if member.User.String() == "stemak#2557" {
				s.ChannelMessageSend(chJoraId, ",add-money "+member.Mention()+" 10")
				s.ChannelMessageSend(chJoraId, member.Mention()+", "+fmt.Sprintf(responces[rand.Intn(len(responces))], "Bump", "1000<:AH_AniCoin:579712087224483850>"))
				fmt.Println("testSTEMAK done")
				break
			}
		}
	}

	if m.ChannelID == chBumpId &&
		len(m.Embeds) > 0 {

		if m.Author.ID == uSupId &&
			m.Embeds[0].Title == "Сервер Up" &&
			m.Embeds[0].Footer != nil {

			guild, err := s.Guild(gAhousId)
			errorLog(err, "ERROR get guild failure:")

			for _, member := range guild.Members {
				if member.User.String() == m.Embeds[0].Footer.Text {

					sendAndLog(s, member.User, "S.up")
					break
				}
			}
		}

		matched, err := regexp.Match(`Server bumped by <@\d*>`, []byte(m.Embeds[0].Description))
		errorLog(err, "ERROR make match regular failure:")

		if matched && m.Author.ID == uBumpId {
			id := strings.Split(strings.Split(m.Embeds[0].Description, "<@")[1], ">")[0]
			user, err := s.User(id)
			errorLog(err, "ERROR get user failure:")

			sendAndLog(s, user, "Bump")
		}
	}
}

func sendAndLog(s *discordgo.Session, member *discordgo.User, str string) {
	s.ChannelMessageSend(chJoraId, ",add-money "+member.Mention()+" 1000")
	s.ChannelMessageSend(chBumpId, member.Mention()+", "+fmt.Sprintf(responces[rand.Intn(len(responces))], str, "1000<:AH_AniCoin:579712087224483850>"))
	fmt.Println("Sever "+str+"ed by", member.String())
}

func errorLog(err error, str string) {
	if err != nil {
		fmt.Println(str, err)
		return
	}
}
