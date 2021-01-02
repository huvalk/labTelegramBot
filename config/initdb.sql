create table if not exists students
    (userId INTEGER NOT NULL PRIMARY KEY,
	fullName VARCHAR(200) NOT NULL,
	groupNum VARCHAR(10) NOT NULL,
	nickname VARCHAR(100) NOT NULL,
	chatId INTEGER NOT NULL);

create table if not exists labs
	(studentId INTEGER NOT NULL,
	labNum INTEGER NOT NULL,
	filePath VARCHAR(300) NOT NULL,
    status VARCHAR(64) NOT NULL,
	messageId INTEGER NOT NULL,
    FOREIGN KEY(studentId) REFERENCES students(userId));

CREATE UNIQUE INDEX if not exists idx_labNum_studentId
    ON labs (studentId, labNum);