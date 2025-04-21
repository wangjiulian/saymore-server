package vo

type StudentPackageList struct {
	ID        int64  `json:"id"`
	StudentId int64  `json:"student_id"`
	Name      string `json:"name"`
	SubjectId int    `json:"subject_id"`
	Hours     string `json:"hours"`
	LeftHours string `json:"left_hours"`
}

type StudentPackageDetail struct {
	ID         int64  `json:"id"`
	StudentId  int64  `json:"student_id"`
	Name       string `json:"name"`
	Hours      string `json:"hours"`
	LeftHours  string `json:"left_hours"`
	Change     uint   `json:"change"`
	ChangeType uint   `json:"change_type"`
	Time       int64  `json:"time"`
}
