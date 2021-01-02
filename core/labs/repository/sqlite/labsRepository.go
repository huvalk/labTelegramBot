package sqliteLabs

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	labsInterfaces "labTelegramBot/core/labs/interfaces"
	baseModels "labTelegramBot/core/models"
)

type labsRepository struct {
	db *sql.DB
}

func NewLabsRepository(db *sql.DB) labsInterfaces.LabsRepository {
	return &labsRepository{
		db: db,
	}
}

func (l *labsRepository) Register(user baseModels.Student) error {
	reg := "INSERT INTO students (userId, fullName, groupNum, nickname, chatId) " +
		"VALUES( $1, $2, $3, $4, $5)"
	_, err := l.db.Exec(reg,
		user.UserID, user.FullName, user.GroupName, user.Nickname, user.ChatID)

	return err
}

func (l *labsRepository) UploadLab(lab baseModels.Lab) error {
	upload := "REPLACE INTO labs (studentId, labNum, filePath, status, messageId) " +
		"VALUES( $1, $2, $3, $4, $5)"
	_, err := l.db.Exec(upload,
		lab.StudentID, lab.LabNum, lab.FilePath, lab.Status, lab.MessageID)

	return err
}

func (l *labsRepository) GetUsers(group string) (students []baseModels.Student, err error) {
	users := "SELECT userId, fullName, nickname FROM students " +
		"WHERE groupNum = $1"
	rows, err := l.db.Query(users, group)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		s := baseModels.Student{}
		err := rows.Scan(&s.UserID, &s.FullName, &s.Nickname)
		if err != nil {
			return nil, err
		}
		students = append(students, s)
	}

	return students, err
}

func (l *labsRepository) GetLabs(userId int) (labs []baseModels.Lab, err error) {
	status := "SELECT labNum, status, filePath FROM labs " +
		"WHERE studentId = $1"
	rows, err := l.db.Query(status, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		l := baseModels.Lab{}
		err := rows.Scan(&l.LabNum, &l.Status, &l.FilePath)
		if err != nil {
			return nil, err
		}
		labs = append(labs, l)
	}

	return labs, err
}
