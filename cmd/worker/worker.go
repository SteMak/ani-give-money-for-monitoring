package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/SteMak/ani-give-money-for-monitoring/bankirapi"

	"github.com/bwmarrin/discordgo"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	api = bankirapi.New(os.Getenv("BANKIRAPI_TOKEN"))

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

	for i := 1; i > 0; i++ {
		time.Sleep(15 * time.Minute)
		fmt.Println("WORKS for", 15*i, "minutes")
	}
}

const (
	chReportsID = "635202206358044710"
	chMonitorID = "569252448137510922"
	chTestID    = "635202206358044710"

	usAdminID = "522347439676588032"
	usSiupID  = "464272403766444044"
	usBumpID  = "315926021457051650"

	glHouseID = "464116508252045312"
)

var (
	err error

	api *bankirapi.API

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

	chMonitorWriters = []simplifiedUser{
		simplifiedUser{
			id:     "522347439676588032",
			strify: "stemak#2557",
		},
	}
)

type simplifiedUser struct {
	id     string
	strify string
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.ChannelID == chMonitorID {
		if len(chMonitorWriters) > 30 {
			chMonitorWriters = append(chMonitorWriters[:1], chMonitorWriters[2:]...)
		}

		chMonitorWriters = append(chMonitorWriters, simplifiedUser{
			id:     m.Author.ID,
			strify: m.Author.String(),
		})
	}

	if m.Author.ID == usAdminID &&
		m.ChannelID == chTestID &&
		strings.HasPrefix(m.Content, "test ") {

		test(s, m)
		return
	}

	if m.ChannelID == chMonitorID &&
		len(m.Embeds) > 0 {

		detectBumpSiup(s, m)
		return
	}
}

func detectBumpSiup(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == usSiupID &&
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

	if matched && m.Author.ID == usBumpID {
		onBumpServer(s, m)
		return
	}
}

func onSiupServer(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("FOUND S.up")

	for _, user := range chMonitorWriters {
		if user.strify == m.Embeds[0].Footer.Text {
			fmt.Println("FOUND S.up user", user.id)

			sendAndLog(s, user.id, "S.up", 1000)
			return
		}
	}
}

func onBumpServer(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("FOUND Bump")

	userID := strings.Split(strings.Split(m.Embeds[0].Description, "<@")[1], ">")[0]
	if len(userID) == 0 {
		fmt.Println("ERROR Bump get user ID:", m.Embeds[0].Description)
		return
	}

	if strings.HasPrefix(userID, "!") {
		userID = userID[1:]
	}

	fmt.Println("FOUND Bump user", userID)

	sendAndLog(s, userID, "Bump", 1000)
}

func test(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "test R U TYT?" {
		s.ChannelMessageSend(chTestID, "E IM TYT!")
		return
	}

	if m.Content == "test S.up" {
		onSiupServerTest(s)
		return
	}

	if m.Content == "test Bump" {
		onBumpServerTest(s)
		return
	}
}

func onSiupServerTest(s *discordgo.Session) {
	fmt.Println("FOUND test S.up")

	for _, user := range chMonitorWriters {
		if user.strify == "stemak#2557" {
			fmt.Println("FOUND test S.up user", user.id)

			sendAndLog(s, user.id, "test S.up", 10)
			return
		}
	}
}

func onBumpServerTest(s *discordgo.Session) {
	fmt.Println("FOUND test Bump")

	userID := strings.Split(strings.Split("Server bumped by <@522347439676588032>. Malades!", "<@")[1], ">")[0]
	if len(userID) == 0 {
		fmt.Println("ERROR test Bump get user ID")
		return
	}

	if strings.HasPrefix(userID, "!") {
		userID = userID[1:]
	}
	fmt.Println("FOUND test Bump user", userID)

	sendAndLog(s, userID, "test Bump", 10)
}

func sendAndLog(s *discordgo.Session, userID string, str string, sum int) {
	var (
		chForLog    = chMonitorID
		chForReport = chReportsID
	)

	if strings.HasPrefix(str, "test ") {
		chForLog = chTestID
		chForReport = chTestID
	}

	bal1, err := api.GetBalance(glHouseID, userID)
	if err != nil {
		fmt.Println("ERROR "+str+" getting user balance:", err)
		return
	}

	bal2, err := api.AddToBalance(glHouseID, userID, 0, sum, "for "+str)
	if err != nil {
		fmt.Println("ERROR "+str+" updating user balance:", err)
		return
	}

	if bal2.Bank-bal1.Bank != 1000 {
		_, err = s.ChannelMessageSend(chForReport, "Кажись, что-то пошло не так... <@"+userID+"> сделал "+str+", но денег ему не дали(")
		if err != nil {
			fmt.Println("ERROR "+str+" sending wrong report message:", err)
		}
		_, err = s.ChannelMessageSend(chForLog, "<@"+userID+">, у нас снова что-то сломалось, но не волнуйтесь - деньги Вам прилетят чуть позже)")
		if err != nil {
			fmt.Println("ERROR "+str+" sending wrong log message:", err)
		}

		return
	}

	_, err = s.ChannelMessageSend(chForReport, strconv.Itoa(sum)+"<:AH_AniCoin:579712087224483850>"+" были выданы <@"+userID+">, за то что он сделал "+str)
	if err != nil {
		fmt.Println("ERROR "+str+" sending right report message:", err)
	}

	_, err = s.ChannelMessageSend(chForLog, "<@"+userID+">, "+fmt.Sprintf(responces[rand.Intn(len(responces))], str, strconv.Itoa(sum)+"<:AH_AniCoin:579712087224483850>"))
	if err != nil {
		fmt.Println("ERROR "+str+" sending right log message:", err)
	}

	fmt.Println("GUILD "+str+" by  ", userID)
}
