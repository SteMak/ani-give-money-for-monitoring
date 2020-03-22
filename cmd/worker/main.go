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
	chJoraID = "635202206358044710"
	chBotID  = "467251523244523522"
	chBumpID = "569252448137510922"

	uStemaID = "522347439676588032"
	uSupID   = "464272403766444044"
	uBumpID  = "315926021457051650"

	gAhousID = "464116508252045312"
)

var (
	err error

	userID string

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
		"Меро уходил спать, но %sнул за Вас и вы получили %s",
		"Нев отвлёк всех разговорами об отсутствии холодильника и вы %sнули сервер. Холодильник дал вам %s",
		"Кемпер поднимал актив, чтобы кемперить сакуру, так что вы тихо %sнули сервер и получили %s",
		"Вы вычислили Маргинала и он рассказал вам секрет %sанья. Получено %s",
		"у Боннуса провис интернет, так что вы беспрепятствеено %sнули сервер и забрали %s",
		"Эксля заснул, что-то бормоча во сне: \"Z-z-z... %s Z-z-z... %s Z-z-z...\"",
		"Эспада зажала Вас в тиски, но вы успели сделать %s и они Вас зауважали и дали %s",
		"Маю-Маю снова уснула в войсе, вы вдохновились её ворочанием и %sнули, заработав %s",
		"ɔıloqoou написал %s перевёрнутыти буквами, поэтому вы неспеша забрали %s",
		"Фузу мирно рисовала в войсе, а Вы сделали %s и собрали %s",
	}
)

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == uStemaID &&
		m.ChannelID == chJoraID &&
		strings.HasPrefix(m.Content, "test ") {

		test(s, m)
		return
	}

	if m.ChannelID == chBumpID &&
		len(m.Embeds) > 0 {

		fmt.Println("FOUND embed in channel of monitoring")

		if m.Author.ID == uSupID &&
			m.Embeds[0].Title == "Сервер Up" &&
			m.Embeds[0].Footer != nil {

			onSiupServer(s, m)
			return
		}

		matched, err := regexp.Match(`Server bumped by <@\d*>`, []byte(m.Embeds[0].Description))
		if err != nil {
			fmt.Println("ERROR Bump make match regular failure:", err)
			return
		}

		if matched && m.Author.ID == uBumpID {
			onBumpServer(s, m)
			return
		}
	}
}

func onSiupServer(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("FOUND S.up message")

	fmt.Println("FOUND S.up user:", m.Embeds[0].Footer.Text)

	guild, err := s.Guild(gAhousID)
	if err != nil {
		fmt.Println("ERROR S.up get guild failure:", err)
		return
	}

	for _, member := range guild.Members {
		if member.User.String() == m.Embeds[0].Footer.Text {

			fmt.Println("FOUND S.up matched member")

			sendAndLog(s, member.User, "S.up")
			break
		}
	}
}

func onBumpServer(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("FOUND Bump message")

	ID := strings.Split(strings.Split(m.Embeds[0].Description, "<@")[1], ">")[0]
	user, err := s.User(ID)
	if err != nil {
		fmt.Println("ERROR Bump get user failure:", err)
		return
	}

	fmt.Println("FOUND Bump user:", user.String())

	sendAndLog(s, user, "Bump")
}

func sendAndLog(s *discordgo.Session, member *discordgo.User, str string) {
	_, err = s.ChannelMessageSend(chJoraID, ",add-money "+member.Mention()+" 1000")
	if err != nil {
		fmt.Println("ERROR "+str+" sending message giving money:", err)
		return
	}

	_, err = s.ChannelMessageSend(chBumpID, member.Mention()+", "+fmt.Sprintf(responces[rand.Intn(len(responces))], str, "1000<:AH_AniCoin:579712087224483850>"))
	if err != nil {
		fmt.Println("ERROR "+str+" sending message notice:", err)
		return
	}

	fmt.Println("GUILD "+str+" by", member.String())
}

func test(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "test R U TYT?" {
		s.ChannelMessageSend(chJoraID, "E IM TYT!")
	}

	if m.Content == "test S.up" {

		fmt.Println("FOUND test S.up message")

		fmt.Println("FOUND test S.up user:", "stemak#2557")

		guild, err := s.Guild(gAhousID)
		if err != nil {
			fmt.Println("ERROR test S.up get guild failure:", err)
			return
		}

		for _, member := range guild.Members {
			if member.User.String() == "stemak#2557" {

				fmt.Println("FOUND test S.up matched member")

				_, err = s.ChannelMessageSend(chJoraID, ",add-money "+member.User.Mention()+" 10")
				if err != nil {
					fmt.Println("ERROR test S.up sending message giving money:", err)
					return
				}

				_, err = s.ChannelMessageSend(chJoraID, member.User.Mention()+", "+fmt.Sprintf(responces[rand.Intn(len(responces))], "test S.up", "10<:AH_AniCoin:579712087224483850>"))
				if err != nil {
					fmt.Println("ERROR test S.up sending message notice:", err)
					return
				}

				fmt.Println("SERVER test S.up by", member.User.String())
				break
			}
		}
	}

	if m.Content == "test Bump" {

		fmt.Println("FOUND test Bump message")

		ID := strings.Split(strings.Split("Server bumped by <@522347439676588032>. Malades!", "<@")[1], ">")[0]
		user, err := s.User(ID)
		if err != nil {
			fmt.Println("ERROR test Bump get user failure:", err)
			return
		}

		fmt.Println("FOUND Bump user:", user.String())

		_, err = s.ChannelMessageSend(chJoraID, ",add-money "+user.Mention()+" 10")
		if err != nil {
			fmt.Println("ERROR test Bump sending message giving money:", err)
			return
		}

		_, err = s.ChannelMessageSend(chJoraID, user.Mention()+", "+fmt.Sprintf(responces[rand.Intn(len(responces))], "test Bump", "10<:AH_AniCoin:579712087224483850>"))
		if err != nil {
			fmt.Println("ERROR test Bump sending message notice:", err)
			return
		}

		fmt.Println("SERVER test Bump by", user.String())
	}
}
