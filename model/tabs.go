package model

type Tab struct {
	ID         int64
	Name       string
	Editor     string
	Output     string
	IsActive   bool
	ActiveDBID *int64
	ActiveDB   *string
}
