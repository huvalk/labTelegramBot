package labsInterfaces

import baseModels "labTelegramBot/core/models"

type LabsRepository interface {
	Register(user baseModels.Student) error
	GetLabs(userId int) (labs []baseModels.Lab, err error)
	GetUsers(group string) (students []baseModels.Student, err error)
	UploadLab(lab baseModels.Lab) error
	SaveMessage(message baseModels.Message) error
	SaveQuestion(message baseModels.Message) error
	GetAllUsersID() ([]int64, error)
}
