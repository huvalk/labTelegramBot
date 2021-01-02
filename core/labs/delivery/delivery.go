package labsDelivery

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/kataras/golog"
	"labTelegramBot/core/labs/interfaces"
)

type labsDelivery struct {
	bot      *tgbotapi.BotAPI
	updates  tgbotapi.UpdatesChannel
	repo     labsInterfaces.LabsRepository
	commands Commands
}

func NewLabsDelivery(bot *tgbotapi.BotAPI, repo labsInterfaces.LabsRepository) (del labsInterfaces.LabsDelivery, err error) {
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	updates, err := bot.GetUpdatesChan(ucfg)
	if err != nil {
		golog.Error("Chanel error: ", err)
		return nil, err
	}

	return &labsDelivery{
		bot:      bot,
		updates:  updates,
		repo:     repo,
		commands: NewCommandList(),
	}, nil
}

func (l *labsDelivery) StartDelivery() error {
	for {
		select {
		case update := <-l.updates:
			if update.Message == nil {
				continue
			}
			if update.Message.Chat.IsGroup() {
				continue
			}

			message := update.Message
			var replay string
			var err error
			switch message.Command() {
			case "register":
				err = l.RegisterUser(message)
				if err != nil {
					golog.Error("register error: ", err)
					replay = err.Error()
				} else {
					replay = fmt.Sprintf("Пользователь добавлен\n ID - %d\n", message.From.ID)
				}

			case "progress":
				replay, err = l.GetStatus(message)
				if err != nil {
					golog.Error("progress error: ", err)
					replay = err.Error()
				}

			case "get_users":
				replay, err = l.GetUsers(message)
				if err != nil {
					golog.Error("get_users error: ", err)
					replay = err.Error()
				}

			case "get_labs":
				replay, err = l.GetLabs(message)
				if err != nil {
					golog.Error("get_labs error: ", err)
					replay = err.Error()
				}

			case "get_all_labs":
				replay, err = l.GetAllLabs(message)
				if err != nil {
					golog.Error("get_labs error: ", err)
					replay = err.Error()
				}

			default:
				err = l.UploadLab(message)
				if err != nil {
					golog.Error("uploadLab error: ", err)
					replay = err.Error()
				} else {
					replay = "Лаба загружена"
				}
			}

			msg := tgbotapi.NewMessage(message.Chat.ID, replay)
			msg.ParseMode = "HTML"
			_, err = l.bot.Send(msg)
			if err != nil {
				golog.Error("sendMessage error: ", err)
			}
		}
	}
}
