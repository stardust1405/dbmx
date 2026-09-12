package app

import (
	"bufio"
	"context"
	"dbmx/model"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	wruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

const (
	ExportFormatCSV  = "csv"
	ExportFormatJSON = "json"
)

// exportQueryTimeout bounds a full-table export. It is far longer than the 30s
// a query in the editor gets: an export is a deliberate act on a table the user
// already knows the size of, and cutting it short would leave a truncated file
// behind with nothing to say it was truncated.
const exportQueryTimeout = 10 * time.Minute

// Export writes result sets to disk. It is bound into the frontend, so every
// exported method here is callable from TypeScript; the save dialog needs the
// Wails context, which only arrives at startup, hence SetExportContext.
type Export struct {
	ctx  context.Context
	conn *Connections
}

func NewExport(conn *Connections) *Export {
	return &Export{conn: conn}
}

// SetExportContext hands Export the Wails context, without which no save dialog
// can open. It is called once from the app's startup hook. It is a function
// rather than a method because every exported method on a bound struct becomes
// part of the frontend's API, and this one is for main to call, not TypeScript.
func SetExportContext(e *Export, ctx context.Context) {
	e.ctx = ctx
}

// exportValue is one cell on its way to a file. NULL is kept apart from the
// empty string because the two formats spell it differently -- an empty field
// in CSV, a literal null in JSON -- and because a text "NULL" would otherwise
// be indistinguishable from the real thing.
type exportValue struct {
	Text   string
	IsNull bool
}

// ExportRows writes rows the grid already holds. Both grids call it for what is
// on screen: the editor's whole result set, and the table view's current page.
func (e *Export) ExportRows(format, suggestedName string, columns []string, rows [][]model.Cell) model.ExportResult {
	if len(columns) == 0 {
		return model.ExportResult{Message: "There is nothing to export."}
	}

	path, err := e.askForPath(format, suggestedName)
	if err != nil {
		return model.ExportResult{Message: err.Error()}
	}
	if path == "" {
		return model.ExportResult{Canceled: true}
	}

	written, err := writeExportFile(path, format, func(w rowWriter) (int64, error) {
		if err := w.WriteHeader(columns); err != nil {
			return 0, err
		}

		values := make([]exportValue, len(columns))
		var count int64
		for _, row := range rows {
			for i := range values {
				if i < len(row) {
					values[i] = cellExportValue(row[i].Value)
				} else {
					// A row shorter than the header is missing those columns
					// outright, which is not the same as holding an empty value.
					values[i] = exportValue{IsNull: true}
				}
			}
			if err := w.WriteRow(values); err != nil {
				return count, err
			}
			count++
		}
		return count, nil
	})
	if err != nil {
		return model.ExportResult{Message: err.Error()}
	}

	return model.ExportResult{OK: true, Path: path, Rows: written}
}

// ExportTable streams every row the table view's filters select, not just the
// page on screen. The clause bar's values are passed in rather than read back
// out of the tabs row, so the export matches the grid even when the user has
// changed a filter since it was last saved.
func (e *Export) ExportTable(tabID int64, format, suggestedName, tableName, selectQuery, where, orderBy, groupBy string) model.ExportResult {
	if strings.TrimSpace(tableName) == "" {
		return model.ExportResult{Message: "A table name is required to export every row."}
	}

	pool, err := e.conn.poolForTab(tabID)
	if err != nil {
		return model.ExportResult{Message: err.Error()}
	}

	path, err := e.askForPath(format, suggestedName)
	if err != nil {
		return model.ExportResult{Message: err.Error()}
	}
	if path == "" {
		return model.ExportResult{Canceled: true}
	}

	ctx, cancel := context.WithTimeout(context.Background(), exportQueryTimeout)
	defer cancel()

	// The same SELECT the grid pages through, without its LIMIT/OFFSET.
	query := buildTableSelect(tableName, selectQuery, where, orderBy, groupBy)

	written, err := writeExportFile(path, format, func(w rowWriter) (int64, error) {
		resultRows, err := pool.Query(ctx, query)
		if err != nil {
			return 0, err
		}
		defer resultRows.Close()

		fields := resultRows.FieldDescriptions()
		columns := make([]string, len(fields))
		for i, field := range fields {
			columns[i] = string(field.Name)
		}
		if err := w.WriteHeader(columns); err != nil {
			return 0, err
		}

		values := make([]exportValue, len(columns))
		var count int64
		for resultRows.Next() {
			raw, err := resultRows.Values()
			if err != nil {
				return count, err
			}
			for i := range values {
				if i < len(raw) {
					values[i] = rawExportValue(raw[i])
				} else {
					values[i] = exportValue{IsNull: true}
				}
			}
			if err := w.WriteRow(values); err != nil {
				return count, err
			}
			count++
		}
		return count, resultRows.Err()
	})
	if err != nil {
		return model.ExportResult{Message: err.Error()}
	}

	return model.ExportResult{OK: true, Path: path, Rows: written}
}

// askForPath puts up the native save dialog. An empty path with no error is the
// user cancelling, which every caller treats as a non-event.
func (e *Export) askForPath(format, suggestedName string) (string, error) {
	if e.ctx == nil {
		return "", errors.New("the export dialog is not available yet")
	}

	var filter wruntime.FileFilter
	switch format {
	case ExportFormatCSV:
		filter = wruntime.FileFilter{DisplayName: "CSV (*.csv)", Pattern: "*.csv"}
	case ExportFormatJSON:
		filter = wruntime.FileFilter{DisplayName: "JSON (*.json)", Pattern: "*.json"}
	default:
		return "", errors.Errorf("unsupported export format %q", format)
	}

	return wruntime.SaveFileDialog(e.ctx, wruntime.SaveDialogOptions{
		Title:                "Export results",
		DefaultFilename:      exportFileName(suggestedName, format),
		Filters:              []wruntime.FileFilter{filter},
		CanCreateDirectories: true,
	})
}

// exportFileName suggests a name the user can accept as-is. The timestamp keeps
// a second export of the same table from silently overwriting the first.
func exportFileName(base, format string) string {
	cleaned := strings.Map(func(r rune) rune {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|':
			return '-'
		}
		return r
	}, strings.TrimSpace(base))

	if strings.TrimSpace(cleaned) == "" {
		cleaned = "export"
	}
	return fmt.Sprintf("%s-%s.%s", cleaned, time.Now().Format("20060102-150405"), format)
}

// writeExportFile owns the file. It fills a temporary file alongside the chosen
// path and only renames it into place once the whole result set is written, so
// a query that fails half way through leaves neither a truncated export at that
// path nor -- where the user picked an existing file to overwrite -- a file
// destroyed by an export that never completed.
func writeExportFile(path, format string, fill func(rowWriter) (int64, error)) (int64, error) {
	// CreateTemp opens at 0600, which the finished export keeps. A file holding
	// whatever the query returned is better off readable by its owner alone.
	file, err := os.CreateTemp(filepath.Dir(path), filepath.Base(path)+".*.part")
	if err != nil {
		return 0, err
	}
	temp := file.Name()

	// Anything that leaves before the rename takes the partial file with it.
	committed := false
	defer func() {
		if !committed {
			os.Remove(temp)
		}
	}()

	writer, err := newRowWriter(format, file)
	if err != nil {
		file.Close()
		return 0, err
	}

	count, err := fill(writer)
	if err == nil {
		err = writer.Close()
	}
	if closeErr := file.Close(); err == nil {
		err = closeErr
	}
	if err != nil {
		return 0, err
	}

	// Same directory, so this is a rename rather than a copy: the file at path
	// is either the previous one or the finished export, never something between.
	if err := os.Rename(temp, path); err != nil {
		return 0, err
	}
	committed = true

	return count, nil
}

// rowWriter serialises one result set. Both implementations take values as
// text-or-null pairs, so a row streamed from postgres and a row the grid handed
// over are written by the same code.
type rowWriter interface {
	WriteHeader(columns []string) error
	WriteRow(values []exportValue) error
	Close() error
}

func newRowWriter(format string, file *os.File) (rowWriter, error) {
	switch format {
	case ExportFormatCSV:
		return &csvRowWriter{w: csv.NewWriter(file)}, nil
	case ExportFormatJSON:
		return &jsonRowWriter{w: bufio.NewWriter(file)}, nil
	default:
		return nil, errors.Errorf("unsupported export format %q", format)
	}
}

type csvRowWriter struct {
	w   *csv.Writer
	row []string
}

func (c *csvRowWriter) WriteHeader(columns []string) error {
	c.row = make([]string, 0, len(columns))
	return c.w.Write(columns)
}

func (c *csvRowWriter) WriteRow(values []exportValue) error {
	c.row = c.row[:0]
	for _, value := range values {
		// A NULL and an empty string both land as an empty field: CSV has no
		// way to tell them apart, and inventing a sentinel would make the file
		// unreadable by anything that did not know about it.
		c.row = append(c.row, value.Text)
	}
	return c.w.Write(c.row)
}

func (c *csvRowWriter) Close() error {
	c.w.Flush()
	return c.w.Error()
}

type jsonRowWriter struct {
	w *bufio.Writer
	// keys are pre-encoded and de-duplicated once, so a result set with two
	// columns of the same name does not write one over the other.
	keys  [][]byte
	count int64
}

func (j *jsonRowWriter) WriteHeader(columns []string) error {
	j.keys = make([][]byte, len(columns))
	for i, name := range uniqueKeys(columns) {
		encoded, err := json.Marshal(name)
		if err != nil {
			return err
		}
		j.keys[i] = encoded
	}
	// bufio.Writer latches its first error and returns it from Flush, so the
	// writes below are left unchecked and Close reports for all of them.
	j.w.WriteString("[")
	return nil
}

func (j *jsonRowWriter) WriteRow(values []exportValue) error {
	if j.count > 0 {
		j.w.WriteString(",")
	}
	j.w.WriteString("\n  {")

	for i, value := range values {
		if i >= len(j.keys) {
			break
		}
		if i > 0 {
			j.w.WriteString(", ")
		}
		j.w.Write(j.keys[i])
		j.w.WriteString(": ")
		if value.IsNull {
			j.w.WriteString("null")
			continue
		}
		encoded, err := json.Marshal(value.Text)
		if err != nil {
			return err
		}
		j.w.Write(encoded)
	}

	j.w.WriteString("}")
	j.count++
	return nil
}

func (j *jsonRowWriter) Close() error {
	if j.count > 0 {
		j.w.WriteString("\n")
	}
	j.w.WriteString("]\n")
	return j.w.Flush()
}

// uniqueKeys makes every column name usable as a distinct JSON key. A result
// set can repeat a name (a join projecting two id columns) or carry none at all
// (a bare expression), neither of which an object can represent as it stands.
func uniqueKeys(columns []string) []string {
	counts := make(map[string]int, len(columns))
	keys := make([]string, len(columns))

	for i, name := range columns {
		if strings.TrimSpace(name) == "" {
			name = fmt.Sprintf("column_%d", i+1)
		}
		key := name
		for counts[key] > 0 {
			counts[name]++
			key = fmt.Sprintf("%s_%d", name, counts[name])
		}
		counts[key]++
		keys[i] = key
	}

	return keys
}

// cellExportValue recovers what it can from a cell the grid already holds. The
// query paths render a SQL NULL as the text "NULL" and an empty string as
// "EMPTY" long before the row reaches the frontend, so those sentinels are the
// only signal left -- which means a column whose real content is the word NULL
// exports as a null too. Rows streamed straight from postgres do not go through
// here and keep the distinction intact.
func cellExportValue(value string) exportValue {
	switch value {
	case "NULL":
		return exportValue{IsNull: true}
	case "EMPTY":
		return exportValue{}
	}
	return exportValue{Text: value}
}

// rawExportValue renders a postgres value the way the grid's own query paths
// render it, so an exported file and the rows on screen read identically.
func rawExportValue(value any) exportValue {
	switch v := value.(type) {
	case nil:
		return exportValue{IsNull: true}
	case []byte:
		return exportValue{Text: string(v)}
	case time.Time:
		return exportValue{Text: v.Format(time.RFC3339)}
	case [16]uint8:
		return exportValue{Text: uuid.UUID(v).String()}
	case string:
		return exportValue{Text: v}
	default:
		return exportValue{Text: fmt.Sprintf("%v", v)}
	}
}
