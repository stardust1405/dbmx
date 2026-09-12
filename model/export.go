package model

// ExportResult is what an export reports back to the grid that asked for it.
// A cancelled save dialog is not a failure, so it gets its own flag rather than
// an OK of false with a message nobody wants to see in a toast.
type ExportResult struct {
	OK       bool `json:"ok"`
	Canceled bool `json:"canceled"`
	// Path is where the file was written, for the confirmation toast.
	Path string `json:"path"`
	// Rows counts the data rows written, not counting the CSV header.
	Rows    int64  `json:"rows"`
	Message string `json:"message"`
}
