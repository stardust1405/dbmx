package app

import (
	"dbmx/model"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// The export writers turn a result set into bytes another program has to read
// back, so what matters is the exact text on disk: how a comma inside a value is
// quoted, how a NULL is spelled in each format, and whether a result set with
// two columns of the same name still round-trips. None of that is visible to the
// compiler, and none of it needs a database.

// writeRows runs the same fill loop ExportRows uses, without the save dialog.
func writeRows(t *testing.T, format string, columns []string, rows [][]model.Cell) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "out."+format)
	count, err := writeExportFile(path, format, func(w rowWriter) (int64, error) {
		if err := w.WriteHeader(columns); err != nil {
			return 0, err
		}
		values := make([]exportValue, len(columns))
		var n int64
		for _, row := range rows {
			for i := range values {
				if i < len(row) {
					values[i] = cellExportValue(row[i].Value)
				} else {
					values[i] = exportValue{IsNull: true}
				}
			}
			if err := w.WriteRow(values); err != nil {
				return n, err
			}
			n++
		}
		return n, nil
	})
	if err != nil {
		t.Fatalf("writeExportFile: %v", err)
	}
	if count != int64(len(rows)) {
		t.Fatalf("wrote %d rows, want %d", count, len(rows))
	}

	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read back: %v", err)
	}
	return string(contents)
}

func cells(columns []string, values ...string) []model.Cell {
	row := make([]model.Cell, len(values))
	for i, value := range values {
		row[i] = model.Cell{Column: columns[i], Value: value}
	}
	return row
}

func TestExportCSVQuotesAndNulls(t *testing.T) {
	columns := []string{"id", "name", "notes"}
	got := writeRows(t, ExportFormatCSV, columns, [][]model.Cell{
		cells(columns, "1", "Ada, Countess", "NULL"),
		cells(columns, "2", "quote \" inside", "EMPTY"),
		cells(columns, "3", "line\nbreak", "plain"),
	})

	// Rows end with LF, encoding/csv's default. A newline inside a quoted value
	// is left exactly as it was, which switching to CRLF would not do.
	want := "id,name,notes\n" +
		"1,\"Ada, Countess\",\n" +
		"2,\"quote \"\" inside\",\n" +
		"3,\"line\nbreak\",plain\n"
	if got != want {
		t.Errorf("csv mismatch\n got: %q\nwant: %q", got, want)
	}
}

func TestExportJSONNullsAndEscaping(t *testing.T) {
	columns := []string{"id", "name"}
	got := writeRows(t, ExportFormatJSON, columns, [][]model.Cell{
		cells(columns, "1", "NULL"),
		cells(columns, "2", "EMPTY"),
		cells(columns, "3", "tab\there"),
	})

	want := "[\n" +
		"  {\"id\": \"1\", \"name\": null},\n" +
		"  {\"id\": \"2\", \"name\": \"\"},\n" +
		"  {\"id\": \"3\", \"name\": \"tab\\there\"}\n" +
		"]\n"
	if got != want {
		t.Errorf("json mismatch\n got: %q\nwant: %q", got, want)
	}
}

// A result set with no rows still has to be a readable file, not an empty one.
func TestExportEmptyResultSets(t *testing.T) {
	columns := []string{"id"}

	if got, want := writeRows(t, ExportFormatCSV, columns, nil), "id\n"; got != want {
		t.Errorf("empty csv = %q, want %q", got, want)
	}
	if got, want := writeRows(t, ExportFormatJSON, columns, nil), "[]\n"; got != want {
		t.Errorf("empty json = %q, want %q", got, want)
	}
}

// A join projecting two id columns is ordinary, and an object cannot hold the
// same key twice -- without this the second column would overwrite the first.
func TestExportJSONDeduplicatesColumnNames(t *testing.T) {
	columns := []string{"id", "id", ""}
	got := writeRows(t, ExportFormatJSON, columns, [][]model.Cell{
		cells(columns, "left", "right", "unnamed"),
	})

	want := "[\n  {\"id\": \"left\", \"id_2\": \"right\", \"column_3\": \"unnamed\"}\n]\n"
	if got != want {
		t.Errorf("json mismatch\n got: %q\nwant: %q", got, want)
	}
}

// A row shorter than the header is missing those columns outright.
func TestExportShortRowPadsWithNull(t *testing.T) {
	columns := []string{"a", "b", "c"}
	got := writeRows(t, ExportFormatJSON, columns, [][]model.Cell{
		{{Column: "a", Value: "1"}},
	})

	want := "[\n  {\"a\": \"1\", \"b\": null, \"c\": null}\n]\n"
	if got != want {
		t.Errorf("json mismatch\n got: %q\nwant: %q", got, want)
	}
}

// A failure part-way through must not leave a truncated file behind at a path
// the user will later mistake for a complete export.
func TestExportRemovesFileWhenFillFails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.csv")

	_, err := writeExportFile(path, ExportFormatCSV, func(w rowWriter) (int64, error) {
		if err := w.WriteHeader([]string{"id"}); err != nil {
			return 0, err
		}
		if err := w.WriteRow([]exportValue{{Text: "1"}}); err != nil {
			return 0, err
		}
		return 1, os.ErrClosed
	})
	if err == nil {
		t.Fatal("expected the fill error to be reported")
	}
	if _, statErr := os.Stat(path); !os.IsNotExist(statErr) {
		t.Errorf("a failed export left %s behind", path)
	}
}

// rawExportValue is what keeps a streamed export reading the same as the grid.
func TestRawExportValueMatchesGridRendering(t *testing.T) {
	stamp := time.Date(2026, 9, 13, 10, 30, 0, 0, time.UTC)

	tests := []struct {
		name   string
		value  any
		want   string
		isNull bool
	}{
		{"nil is null", nil, "", true},
		{"bytes become text", []byte("bytes"), "bytes", false},
		{"time is rfc3339", stamp, "2026-09-13T10:30:00Z", false},
		{"empty string stays empty", "", "", false},
		{"numbers are formatted", int64(42), "42", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := rawExportValue(test.value)
			if got.Text != test.want || got.IsNull != test.isNull {
				t.Errorf("rawExportValue(%v) = %+v, want {Text:%q IsNull:%v}",
					test.value, got, test.want, test.isNull)
			}
		})
	}
}

// The export runs the same SELECT the grid pages through; if they drift, an
// export silently stops matching the filters on screen.
func TestBuildTableSelect(t *testing.T) {
	tests := []struct {
		name                                        string
		table, selectQuery, where, orderBy, groupBy string
		want                                        string
	}{
		{
			name:  "defaults to every column ordered by the first",
			table: "users",
			want:  `SELECT * FROM "users" ORDER BY 1`,
		},
		{
			name:        "carries every clause in SQL order",
			table:       "users",
			selectQuery: "id, email",
			where:       "active = true",
			groupBy:     "id, email",
			orderBy:     "email DESC",
			want:        `SELECT id, email FROM "users" WHERE active = true GROUP BY id, email ORDER BY email DESC`,
		},
		{
			name:        "blank clauses are dropped, not rendered empty",
			table:       "users",
			selectQuery: "  ",
			where:       "   ",
			want:        `SELECT * FROM "users" ORDER BY 1`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := buildTableSelect(test.table, test.selectQuery, test.where, test.orderBy, test.groupBy)
			if got != test.want {
				t.Errorf("buildTableSelect()\n got: %s\nwant: %s", got, test.want)
			}
		})
	}
}

func TestExportFileNameSanitisesPathSeparators(t *testing.T) {
	got := exportFileName("public/users", ExportFormatCSV)
	if want := "public-users-"; len(got) < len(want) || got[:len(want)] != want {
		t.Errorf("exportFileName = %q, want it to start with %q", got, want)
	}
	if got[len(got)-4:] != ".csv" {
		t.Errorf("exportFileName = %q, want a .csv extension", got)
	}

	if got := exportFileName("   ", ExportFormatJSON); got[:7] != "export-" {
		t.Errorf("exportFileName(blank) = %q, want it to fall back to export-", got)
	}
}
