package format

import (
	"bytes"
	"text/template"
)

// RenderTemplate renders a template string with the given data.
// Returns "[tpl-err]" on error to indicate template issues in the status line.
func RenderTemplate(tmplStr string, data interface{}) string {
	if tmplStr == "" {
		return ""
	}

	tmpl, err := template.New("").Parse(tmplStr)
	if err != nil {
		return "[tpl-err]" // Invalid template syntax
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, data); err != nil {
		return "[tpl-err]" // Template execution failed
	}

	return buf.String()
}
