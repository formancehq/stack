package modules

import "text/template"

type PullConfiguration struct {
	// ModuleURLTpl is the go template used to build url to fetch the module logs
	ModuleURLTpl *template.Template

	// PullPageSize is the page size used to pull modules
	PullPageSize int
}
