create table if not exists students
    (id INTEGER NOT NULL PRIMARY KEY,
	fullName VARCHAR(200) NOT NULL,
	groupNum VARCHAR(10) NOT NULL,
	nickname VARCHAR(100) NOT NULL,
	chatId INTEGER NOT NULL);

create table if not exists labs
	(id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
    student_id INTEGER NOT NULL,
	labNum INTEGER NOT NULL,
	filePath VARCHAR(300) NOT NULL,
    status VARCHAR(64) NOT NULL,
	messageId INTEGER NOT NULL,
    FOREIGN KEY(student_id) REFERENCES students(id));

create table if not exists messages
    (user_id INTEGER NOT NULL,
     chatId INTEGER NOT NULL,
     messageId INTEGER NOT NULL,
     text TEXT NOT NULL,
     addition TEXT NOT NULL,
     sent TIMESTAMP);

create table if not exists questions
(
    user_id   INTEGER not null,
    chatId    INTEGER not null,
    messageId INTEGER not null,
    text      TEXT    not null,
    addition  TEXT    not null,
    sent TIMESTAMP
);

CREATE UNIQUE INDEX if not exists idx_labNum_studentId
    ON labs (student_id, labNum);