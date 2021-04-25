package evil

import (
	"embed"
	"net/http"
)

//go:embed *.js
var Static embed.FS

var StaticHandler = http.FileServer(http.FS(Static))
