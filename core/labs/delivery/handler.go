package labsDelivery

import (
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/kataras/golog"
	"github.com/pkg/errors"
	baseModels "labTelegramBot/core/models"
	"strings"
)

func (l *labsDelivery) RegisterUser(message *tgbotapi.Message) (err error) {
	student := baseModels.Student{
		UserID:    message.From.ID,
		FullName:  "",
		GroupName: "",
		Nickname:  message.From.UserName,
		ChatID:    message.Chat.ID,
	}

	student.FullName, student.GroupName, err =
		l.commands.ExtractGroupAndName(message.CommandArguments())
	if err != nil {
		return err
	}

	err = l.repo.Register(student)

	return err
}

func (l *labsDelivery) GetStatus(message *tgbotapi.Message) (result string, err error) {
	student := baseModels.Student{
		UserID:    message.From.ID,
		FullName:  "",
		GroupName: "",
		Nickname:  message.From.UserName,
		ChatID:    message.Chat.ID,
	}

	labs, err := l.repo.GetLabs(student.UserID)
	if err != nil {
		return "", err
	}
	result = "Лабы:\n"
	for _, l := range labs {
		result += fmt.Sprintf("Лаба %d - %s\n", l.LabNum, l.Status)
	}

	return result, err
}

func (l *labsDelivery) UploadLab(message *tgbotapi.Message) (err error) {
	doc := message.Document
	if doc == nil {
		return errors.Errorf("Повторите снова, не прикреплен pdf с лабой")
	} else if !strings.Contains(doc.FileName, ".pdf") {
		return errors.Errorf("Повторите снова, файл должен быть pdf")
	}

	fileURL, err := l.bot.GetFileDirectURL(doc.FileID)
	if err != nil {
		golog.Error("getFileUrl error: ", err)
		return errors.Errorf("Ошибка загрузки файла")
	}

	lab := baseModels.Lab{
		StudentID: message.From.ID,
		LabNum:    0,
		FilePath:  fileURL,
		Status:    "сдано",
		MessageID: message.MessageID,
	}

	lab.LabNum, err = l.commands.ExtractLab(message.Caption)
	if err != nil {
		return errors.Errorf("Повторите снова, не указан номер лабы")
	}
	err = l.repo.UploadLab(lab)

	return err
}

func (l *labsDelivery) GetUsers(message *tgbotapi.Message) (result string, err error) {
	group, err := l.commands.ExtractGroup(message.CommandArguments())
	if err != nil {
		return "", errors.Errorf("Повторите снова, не указан номер группы")
	}
	students, err := l.repo.GetUsers(group)

	result = "Группа " + group + ":\n"
	for _, s := range students {
		result += fmt.Sprintf("%s - %s - %d\n", s.FullName, s.Nickname, s.UserID)
	}

	return result, err
}

func (l *labsDelivery) GetLabs(message *tgbotapi.Message) (result string, err error) {
	userID, err := l.commands.ExtractStudentID(message.CommandArguments())
	if err != nil {
		return "", errors.Errorf("Повторите снова, не указан номер студента")
	}

	labs, err := l.repo.GetLabs(userID)
	if err != nil {
		return "", err
	}
	result = "Лабы:\n"
	for _, l := range labs {
		result += fmt.Sprintf("<b>Лаба %d - %s\n%s\n</b>", l.LabNum, l.Status, l.FilePath)
	}

	return result, err
}

func (l *labsDelivery) GetAllLabs(message *tgbotapi.Message) (result string, err error) {
	group, err := l.commands.ExtractGroup(message.CommandArguments())
	if err != nil {
		return "", errors.Errorf("Повторите снова, не указан номер группы")
	}
	students, err := l.repo.GetUsers(group)

	result = "<b>Группа " + group + ":\n</b>"
	for _, s := range students {
		labs, err := l.repo.GetLabs(s.UserID)
		if err != nil {
			return "", err
		}
		result += fmt.Sprintf("<b>%s - %s - %d\n</b>", s.FullName, s.Nickname, s.UserID)
		for _, l := range labs {
			result += fmt.Sprintf("Лаба %d - %s\n%s\n", l.LabNum, l.Status, l.FilePath)
		}
	}

	return result, err
}

func (l *labsDelivery) SaveMessage(message *tgbotapi.Message) (err error) {
	messageToSave := baseModels.Message{
		StudentID:  message.From.ID,
		MessageID:  message.MessageID,
		ChatID:     message.Chat.ID,
		Message:    message.CommandArguments() + "\n" + message.Text + "\n" + message.Caption,
		Additional: message.From.FirstName + "\n" + message.From.LastName + "\n" + message.From.UserName,
	}

	err = l.repo.SaveMessage(messageToSave)

	return err
}

func (l *labsDelivery) SaveQuestion(message *tgbotapi.Message) (err error) {
	messageToSave := baseModels.Message{
		StudentID:  message.From.ID,
		MessageID:  message.MessageID,
		ChatID:     message.Chat.ID,
		Message:    message.CommandArguments(),
		Additional: message.From.FirstName + "\n" + message.From.LastName + "\n" + message.From.UserName,
	}

	err = l.repo.SaveQuestion(messageToSave)

	return err
}