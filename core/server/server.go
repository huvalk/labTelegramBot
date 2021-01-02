package server

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/kataras/golog"
	labsDelivery "labTelegramBot/core/labs/delivery"
	"labTelegramBot/core/labs/interfaces"
	"labTelegramBot/core/labs/repository/sqlite"
	"labTelegramBot/core/utils/db"
	"os"
)

type App struct {
	del labsInterfaces.LabsDelivery
}

func NewApp() (app *App, err error) {
	db, err := dbUtil.NewDataBase()
	if err != nil {
		golog.Error("DB error: ", err.Error())
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		golog.Error("DB: ", err.Error())
		return nil, err
	}

	bot, err := newBot()
	labsRepo := sqliteLabs.NewLabsRepository(db)
	del, err := labsDelivery.NewLabsDelivery(bot, labsRepo)

	return &App{
		del: del,
	}, nil
}

func newBot() (bot *tgbotapi.BotAPI, err error) {
	botToken := os.Getenv("BOT_TOKEN")
	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		golog.Error("Bot error: ", err)

		return nil, err
	}

	bot.Debug = false
	golog.Infof("Authorized on account %s", bot.Self.UserName)
	return bot, nil
}

func (a *App) StartApp() error {
	return a.del.StartDelivery()
}
