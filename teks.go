package teks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"
	"text/template"
)

// TableFormatKey is the identifier used for table
const TableFormatKey = "table"

// basicFuncs are used for common data printing
var basicFuncs = template.FuncMap{
	"json": func(v interface{}) string {
		buf := &bytes.Buffer{}
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(false)
		_ = encoder.Encode(v)
		return strings.TrimSpace(buf.String())
	},
	"jsonPretty": func(v interface{}) string {
		buf := &bytes.Buffer{}
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(false)
		encoder.SetIndent("", "  ")
		_ = encoder.Encode(v)
		return strings.TrimSpace(buf.String())
	},
	"split": strings.Split,
	"upper": strings.ToUpper,
	"lower": strings.ToLower,
	"title": strings.Title,
	"join":  strings.Join,
}

// HeaderFuncs are used to format headers in a table
// Some of functions in basicFuncs are overridden
var HeaderFuncs = template.FuncMap{
	"json":       func(s string) string { return s },
	"jsonPretty": func(s string) string { return s },
	"join":       func(s string) string { return s },
}

// Format is an alias for a string used for formatting
type Format string

// IsTable returns true if format string is prefixed with table
func (f Format) IsTable() bool {
	return strings.HasPrefix(string(f), TableFormatKey)
}

// Context keeps data about a format operation
type Context struct {
	// Output is used to write the output
	Output io.Writer
	// Format is used to keep format string
	Format Format

	// internal usage
	finalFormat string
	buffer      *bytes.Buffer
}

// NewContext creates a context with initialized fields
func NewContext(output io.Writer, format string) *Context {
	return &Context{Output: output, Format: Format(format), buffer: &bytes.Buffer{}}
}

// preFormat will clean format string
func (ctx *Context) preFormat() {
	format := string(ctx.Format)

	if ctx.Format.IsTable() {
		// if table is found skip it and take the rest
		format = format[len(TableFormatKey):]
	}

	format = strings.TrimSpace(format)
	// this is done to avoid treating \t \n as template strings. This replaces them as special characters
	replacer := strings.NewReplacer(`\t`, "\t", `\n`, "\n")
	format = replacer.Replace(format)
	ctx.finalFormat = format
}

// parseTemplate will create a new template with basic functions
func (ctx Context) parseTemplate() (*template.Template, error) {
	tmpl, err := NewBasicFormatter("").Parse(ctx.finalFormat)
	if err != nil {
		return tmpl, fmt.Errorf("Template parsing error: %v\n", err)
	}
	return tmpl, nil
}

// postFormat will output to writer
func (ctx *Context) postFormat(template *template.Template, headers interface{}) {
	if ctx.Format.IsTable() {
		// create a tab writer using Output
		w := tabwriter.NewWriter(ctx.Output, 20, 1, 3, ' ', 0)
		// print headers
		_ = template.Funcs(HeaderFuncs).Execute(w, headers)
		_, _ = w.Write([]byte{'\n'})
		// write buffer to the w
		// in this case anything in buffer will be rendered by tab writer to the Output
		// buffer contains data to be written
		_, _ = ctx.buffer.WriteTo(w)
		// flush will perform actual write to the writer
		_ = w.Flush()
	} else {
		// just write it as normal
		_, _ = ctx.buffer.WriteTo(ctx.Output)
	}
}

// Renderer is used to render a particular resource using templates
type Renderer func(io.Writer, *template.Template) error

// Write writes data using r and headers
func (ctx *Context) Write(r Renderer, headers interface{}) error {
	// prepare formatting
	ctx.preFormat()
	// parse template
	tmpl, err := ctx.parseTemplate()
	if err != nil {
		return err
	}
	// using renderer provided render collection
	// Note: See the renderer implementation in cmd/apis.go for more
	if err = r(ctx.buffer, tmpl); err != nil {
		return err
	}
	// write results to writer
	ctx.postFormat(tmpl, headers)
	return nil
}

// NewBasicFormatter creates a new template engine with name
func NewBasicFormatter(name string) *template.Template {
	tmpl := template.New(name).Funcs(basicFuncs)
	return tmpl
}
