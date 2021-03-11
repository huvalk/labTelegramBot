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
			message := update.Message
			if message == nil {
				message = update.EditedMessage
			}
			if message == nil {
				msg := tgbotapi.NewMessage(message.Chat.ID, "Что-то пошло не так, отправьте сообщение заново")
				msg.ParseMode = "HTML"
				_, err := l.bot.Send(msg)
				if err != nil {
					golog.Error("sendMessage error: ", err)
				}
				continue
			}

			if message.Chat.IsGroup() {
				continue
			}

			err := l.SaveMessage(message)
			if err != nil {
				golog.Error("Cant save log message: ", err)
			}

			var replay string
			switch message.Command() {
			case "help":
				replay = "/register ФИО ИУ7-4X - для регистрации (обязательно)\n" +
					"*номер лабы* + pdf одним  сообщением - отправить лабу\n" +
					"/progress - увидеть статус своих лабы\n" +
					"/question ВОПРОС - задать вопрос\n\n" +
					"ВАЖНО Бот не хранит состояние, одно действие - одно сообщение"
			case "start":
				replay = "/register ФИО ИУ7-4X - для регистрации (обязательно)\n" +
					"*номер лабы* + pdf одним  сообщением - отправить лабу\n" +
					"/progress - увидеть статус своих лабы\n" +
					"/question ВОПРОС - задать вопрос\n\n" +
					"ВАЖНО Бот не хранит состояние, одно действие - одно сообщение"
			case "question":
				err = l.SaveQuestion(message)
				if err != nil {
					golog.Error("question error: ", err)
					replay = err.Error()
				} else {
					replay = "Вопрос задан"
				}
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
