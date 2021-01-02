package labsDelivery

import (
	"errors"
	"regexp"
	"strconv"
)

type Commands struct {
	fullName *regexp.Regexp
	group    *regexp.Regexp
	lab    *regexp.Regexp
	studentID    *regexp.Regexp
}

func NewCommandList() Commands {
	return Commands{
		fullName: compileRegexp(`(?:[А-Яа-я]+\s){1,3}[А-Яа-я]+(?:\sИУ)`),
		group:    compileRegexp(`ИУ\d-\d{2}`),
		lab:    compileRegexp(`\d+`),
		studentID:    compileRegexp(`\d+`),
	}
}

func compileRegexp(s string) *regexp.Regexp {
	r, _ := regexp.Compile(s)
	return r
}

func (c *Commands) ExtractGroup(t string) (group string, err error) {
	groupArr := c.group.FindStringSubmatch(t)

	if len(groupArr) != 1 {
		return "", errors.New("Неверные аргументы")
	}

	group = groupArr[0]
	if len(group) > 9 {
		return "", errors.New("Слишком длинный номер группы")
	}

	return group, nil
}

func (c *Commands) ExtractGroupAndName(t string) (fullName string, group string, err error) {
	fullNameArr := c.fullName.FindStringSubmatch(t)
	group, err = c.ExtractGroup(t)

	if err != nil || len(fullNameArr) != 1 {
		return "", "", errors.New("Неверные аргументы")
	}

	fullName = fullNameArr[0]
	if len(fullName) > 199 {
		return "", "", errors.New("Слишком длинное имя")
	}

	return fullName[:len(fullName)-5], group, nil
}

func (c *Commands) ExtractLab(t string) (labNum int, err error) {
	labArr := c.lab.FindStringSubmatch(t)

	if len(labArr) != 1 {
		return 0, errors.New("Неверные аргументы")
	}

	labNum, err = strconv.Atoi(labArr[0])
	if err != nil {
		return 0, errors.New("Неверные аргументы")
	}

	return labNum, nil
}

func (c *Commands) ExtractStudentID(t string) (studentID int, err error) {
	labArr := c.lab.FindStringSubmatch(t)

	if len(labArr) != 1 {
		return 0, errors.New("Неверные аргументы")
	}

	studentID, err = strconv.Atoi(labArr[0])
	if err != nil {
		return 0, errors.New("Неверные аргументы")
	}

	return studentID, nil
}