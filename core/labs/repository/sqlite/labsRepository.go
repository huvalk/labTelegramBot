package sqliteLabs

import (
	"database/sql"
	"fmt"
	"github.com/kataras/golog"
	_ "github.com/mattn/go-sqlite3"
	"io"
	labsInterfaces "labTelegramBot/core/labs/interfaces"
	baseModels "labTelegramBot/core/models"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type labsRepository struct {
	db   *sql.DB
	path string
}

func NewLabsRepository(db *sql.DB) labsInterfaces.LabsRepository {
	p, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	p = path.Dir(p)
	p = path.Join("data", "files")
	if err = os.MkdirAll(p, os.ModePerm); err != nil {
		panic(err)
	}
	return &labsRepository{
		path: p,
		db:   db,
	}
}

func (l *labsRepository) saveFile(group, name, url string) error {
	dirName := filepath.Join(l.path, group, name)

	files, err := os.ReadDir(dirName)
	if err != nil && !os.IsNotExist(err) {
		return err
	} else if os.IsNotExist(err) {
		if err = os.MkdirAll(dirName, os.ModePerm); err != nil {
			panic(err)
		}
	}
	fileName := fmt.Sprintf("%s %s (%d).pdf", group, name, len(files))

	out, err := os.Create(path.Join(dirName, fileName))
	if err != nil {
		return err
	}
	defer out.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		golog.Errorf("Failed to load file \"%s\". Status: %d, url: %s\n", fileName, resp.StatusCode, url)
	}
	_, err = io.Copy(out, resp.Body)
	return err
}

func (l *labsRepository) Register(user baseModels.Student) error {
	reg := "INSERT INTO students (id, fullName, groupNum, nickname, chatId) " +
		"VALUES( $1, $2, $3, $4, $5)"
	_, err := l.db.Exec(reg,
		user.UserID, user.FullName, user.GroupName, user.Nickname, user.ChatID)

	return err
}

func (l *labsRepository) uploadLab(lab baseModels.Lab) error {
	upload := "REPLACE INTO labs (student_id, labNum, filePath, status, messageId) " +
		"VALUES( $1, $2, $3, $4, $5)"
	_, err := l.db.Exec(upload,
		lab.StudentID, lab.LabNum, lab.FilePath, lab.Status, lab.MessageID)

	return err
}

func (l *labsRepository) UploadLab(lab baseModels.Lab) error {
	student := &baseModels.Student{}
	users := "SELECT fullName, groupNum FROM students " +
		"WHERE id = $1"
	row := l.db.QueryRow(users, lab.StudentID)
	if err := row.Scan(&student.FullName, &student.GroupName); err != nil {
		return err
	}

	if err := l.saveFile(student.GroupName, student.FullName, lab.FilePath); err != nil {
		return err
	}
	if err := l.uploadLab(lab); err != nil {
		return err
	}

	return nil
}

func (l *labsRepository) GetUsers(group string) (students []baseModels.Student, err error) {
	users := "SELECT id, fullName, nickname FROM students " +
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
		"WHERE student_id = $1"
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
		l.FilePath = ""
		labs = append(labs, l)
	}

	return labs, err
}

func (l *labsRepository) SaveMessage(message baseModels.Message) (err error) {
	reg := "INSERT INTO messages (user_id, chatId, messageId, text, addition, sent) " +
		"VALUES( $1, $2, $3, $4, $5, $6)"
	_, err = l.db.Exec(reg,
		message.StudentID, message.ChatID, message.MessageID, message.Message, message.Additional, message.Time)

	return err
}

func (l *labsRepository) SaveQuestion(message baseModels.Message) error {
	reg := "INSERT INTO questions (user_id, chatId, messageId, text, addition, sent) " +
		"VALUES( $1, $2, $3, $4, $5, $6)"
	_, err := l.db.Exec(reg,
		message.StudentID, message.ChatID, message.MessageID, message.Message, message.Additional, message.Time)

	return err
}

func (l *labsRepository) GetAllUsersID() (IDs []int64, err error) {
	status := "SELECT chatId FROM students"
	rows, err := l.db.Query(status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var ID int64
		err := rows.Scan(&ID)
		if err != nil {
			return nil, err
		}
		IDs = append(IDs, ID)
	}

	return IDs, err
}
