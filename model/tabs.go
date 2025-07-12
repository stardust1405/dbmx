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
	Type          string

	// Output
	Columns []string `json:"columns"`
	Rows    [][]Cell `json:"rows"`
}

var validTypes = map[string]struct{}{
	"editor": {},
	"table":  {},
}

func IsValidTabType(t string) bool {
	_, ok := validTypes[t]
	return ok
}
