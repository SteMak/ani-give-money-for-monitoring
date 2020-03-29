package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/SteMak/ani-give-money-for-monitoring/workerTools/bankirapi"
	"github.com/SteMak/ani-give-money-for-monitoring/workerTools/config"

	"github.com/bwmarrin/discordgo"
)

var (
	err error

	api *bankirapi.API

	chMonitorWriters []simplifiedUser
)

type simplifiedUser struct {
	id     string
	strify string
}

func main() {
	config.Init()
	
	rand.Seed(time.Now().UnixNano())

	api = bankirapi.New(config.BankirToken)

	fmt.Println("1 WORKER started")

	dg, err := discordgo.New("Bot " + config.Token)
	if err != nil {
		fmt.Println("ERROR creating Discord session:", err)
		return
	}

	fmt.Println("2 Discord session created")

	dg.AddHandler(messageCreate)

	fmt.Println("3 Registred the messageCreate handler")

	err = dg.Open()
	if err != nil {
		fmt.Println("ERROR opening connection:", err)
		return
	}

	fmt.Println("4 Opened a websocket")

	for i := 1; i > 0; i++ {
		time.Sleep(25 * time.Minute)
		fmt.Println("WORKS for", 25*i, "minutes")
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "R U TYT?" && m.ChannelID == config.ChReportsID {
		s.ChannelMessageSend(m.ChannelID, "E IM TYT!")
	}

	if m.ChannelID == config.ChMonitorID {
		if len(chMonitorWriters) >= 30 {
			chMonitorWriters = chMonitorWriters[1:]
		}

		chMonitorWriters = append(chMonitorWriters, simplifiedUser{
			id:     m.Author.ID,
			strify: m.Author.String(),
		})
	}

	if m.ChannelID == config.ChMonitorID && len(m.Embeds) > 0 {
		detectBumpSiup(s, m)
		return
	}
}

func detectBumpSiup(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == config.UsSiupID &&
		m.Embeds[0].Title == "Сервер Up" &&
		m.Embeds[0].Footer != nil {

		onSiupServer(s, m)
		return
	}

	matched, err := regexp.Match(`Server bumped by <@\d+>`, []byte(m.Embeds[0].Description))
	if err != nil {
		fmt.Println("ERROR Bump make match regular failure:", err)
		return
	}

	if matched && m.Author.ID == config.UsBumpID {
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

func sendAndLog(s *discordgo.Session, userID string, str string, sum int) {
	_, err = api.AddToBalance(config.GlHouseID, userID, 0, sum, "for "+str)
	if err != nil {
		fmt.Println("ERROR "+str+" updating user balance:", err)

		_, err = s.ChannelMessageSend(config.ChReportsID, "Кажись, что-то пошло не так... <@"+userID+"> сделал "+str+", но денег ему не дали(")
		if err != nil {
			fmt.Println("ERROR "+str+" sending wrong report message:", err)
		}
		_, err = s.ChannelMessageSend(config.ChMonitorID, "<@"+userID+">, у нас снова что-то сломалось, но не волнуйтесь - деньги Вам прилетят чуть позже)")
		if err != nil {
			fmt.Println("ERROR "+str+" sending wrong log message:", err)
		}

		return
	}

	_, err = s.ChannelMessageSend(config.ChReportsID, strconv.Itoa(sum)+"<:AH_AniCoin:579712087224483850> были выданы <@"+userID+">, за то что он сделал "+str)
	if err != nil {
		fmt.Println("ERROR "+str+" sending right report message:", err)
	}

	_, err = s.ChannelMessageSend(config.ChMonitorID, "<@"+userID+">, "+fmt.Sprintf(config.Responces[rand.Intn(len(config.Responces))], str, strconv.Itoa(sum)+"<:AH_AniCoin:579712087224483850>"))
	if err != nil {
		fmt.Println("ERROR "+str+" sending right log message:", err)
	}

	fmt.Println("GUILD "+str+" by  ", userID)
}
