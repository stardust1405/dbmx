//go:build !production

package database

import _ "embed"

// Dev build embeds app.dev.db (must exist at compile time).
//
//go:embed app.dev.db
var embeddedDB []byte

// Optional: if you want a different runtime filename for dev, set this var.
// Otherwise delete it and use "app.db" in sqlite3.go.
var runtimeDBName = "app.db"
