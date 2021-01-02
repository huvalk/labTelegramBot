package baseModels

type Student struct {
	UserID    int
	FullName  string
	GroupName string
	Nickname  string
	ChatID    int64
}

type Lab struct {
	StudentID int
	LabNum    int
	FilePath  string
	Status    string
	MessageID int
}
