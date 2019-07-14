package config

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/gin-gonic/gin"
	"github.com/tfwio/sekhem/util"
	"github.com/tfwio/sekhem/util/pandoc"
)

type pandocData struct {
	Document string
}

// pandoctemplate expects a structure with `{"Document": <string> }`
var pandoctemplate = template.Must(
	template.New("pandoc").
		Funcs(template.FuncMap{"noDown": noDown}).
		Parse(`{{define "pandoc"}}{{.Document | noDown}}{{end}}`))

// mandocmetatemplate expects a structure with `{"Document": <string> }`
var mandocmetatemplate = template.Must(
	template.New("meta").
		Funcs(template.FuncMap{"noDown": noDown}).
		Parse(`{{define "meta"}}{{.Document | noDown}}{{end}}`))

func (c *Configuration) servePandoc(pandocTemplate string, t *template.Template, g *gin.Context) {

	route := g.Param("path")
	action := g.Param("action")
	if path, err := c.GetFilePath(route, action); err == nil {
		exts, flags := c.Pandoc.Extensions, c.Pandoc.Flags
		if exts == "" {
			exts = "+abbreviations+auto_identifiers+autolink_bare_uris+backtick_code_blocks+bracketed_spans+definition_lists+emoji+escaped_line_breaks+example_lists+fancy_lists+fenced_code_attributes+fenced_divs+footnotes+header_attributes+inline_code_attributes+implicit_figures+implicit_header_references+inline_notes+link_attributes+mmd_title_block+multiline_tables+raw_tex+simple_tables+smart+startnum+strikeout+table_captions+yaml_metadata_block"
		}
		if flags == "" {
			flags = "-N" // numbered header indexes
		}
		var wrap = pandoc.Create(
			util.Abs(c.Pandoc.Executable),
			flags,
			exts,
			util.Abs(pandocTemplate))

		var mByteBuffer bytes.Buffer
		if err := wrap.Do(path, &mByteBuffer, nil, true); err == nil {
			html := pandocData{Document: mByteBuffer.String()}
			if e := t.ExecuteTemplate(g.Writer, "pandoc", html); e != nil {
				fmt.Println("  - pandoc error:", e)
			}
		} else {
			fmt.Println("- NOT pandoc processing", err)
			html := pandocData{Document: "Some kind of error."}
			if e := t.ExecuteTemplate(g.Writer, "pandoc", html); e != nil {
				fmt.Println("Unexpected error:", e)
			}
		}
	}
}
