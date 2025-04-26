package model

type Tab struct {
	ID            int64
	Name          string
	Editor        string
	Output        string
	IsActive      bool
	ActiveDBID    *string
	ActiveDB      *string
	ActiveDBColor *string

	// Output
	Columns []string `json:"columns"`
	Rows    [][]Cell `json:"rows"`
}
