package main

import (
	"html/template"
	"path/filepath"
	"strings"
	"time"
)

// Ctx is used to represent command context.
type Ctx struct {
	Cmd   string
	Input string
	Start time.Time
	Time  time.Time
}

var tmplFuncs = template.FuncMap{
	"toUpper":      strings.ToUpper,
	"toLower":      strings.ToLower,
	"absolutePath": filepath.Abs,
	"basename":     filepath.Base,
	"dirname":      filepath.Dir,
	"ext":          filepath.Ext,
	"noExt":        noExt,
}

func noExt(str string) string {
	return strings.TrimSuffix(str, filepath.Ext(str))
}
