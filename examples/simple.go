package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"text/template"

	"github.com/kasvith/teks"
)

// Person represents a human being or whatever
type Person struct {
	Name    string
	Age     int
	Address string
}

func main() {
	var format string
	// default one will printout Name and Age in tabular format
	flag.StringVar(&format, "format", "table {{.Name}}\t{{.Age}}", "format of output")
	flag.Parse()

	// whatever data you have
	persons := []Person{
		{"Kasun", 24, "Earth"},
		{"Jhon Doe", 34, "Somewhere on earth"},
		{"Spongebob", 30, "Under Sea"},
		{"Harry Potter", 30, "4 Privet Drive, Little Whinging, Surrey"},
	}

	// create new context
	ctx := teks.NewContext(os.Stdout, format)

	// create a renderer function to match signature
	renderer := func(w io.Writer, t *template.Template) error {
		for _, p := range persons {
			if err := t.Execute(w, p); err != nil {
				return err
			}
			_, _ = w.Write([]byte{'\n'})
		}
		return nil
	}

	// headers for table
	tableHeaders := map[string]string{
		"Age":     "Age",
		"Name":    "Name",
		"Address": "Address",
	}

	//override header functions if you want
	//teks.HeaderFuncs = template.FuncMap{
	//	"split": strings.Split,
	//}

	// execute context and write to our output
	if err := ctx.Write(renderer, tableHeaders); err != nil {
		fmt.Println("Error executing template:", err.Error())
	}
}
