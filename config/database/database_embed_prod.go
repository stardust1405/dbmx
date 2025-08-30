//go:build production

package database

import _ "embed"

// Embed app.db that lives NEXT TO this file.
//
//go:embed app.db
var embeddedDB []byte

var runtimeDBName = "app.db"
