# teks : awesome outputs for your commands

[![Build Status](https://travis-ci.com/kasvith/teks.svg?branch=master)](https://travis-ci.com/kasvith/teks) [![Go Report Card](https://goreportcard.com/badge/github.com/kasvith/teks)](https://goreportcard.com/report/github.com/kasvith/teks) [![GoDoc](https://godoc.org/github.com/kasvith/teks?status.svg)](https://godoc.org/github.com/kasvith/teks)

**teks** brings painless output formating for your commands. Docker/Kubernetes provides custom formatting via `go-templates`.
**teks** brings the power into any application providing a smooth intergration as a library. 

> **teks** is hevily inspired by [Docker CLI](https://github.com/docker/cli)

# Install

`teks` is a go package. To use it execute
```
go get github.com/kasvith/teks
```

# Available formatting options

| Name | Usage |
| --- | --- |
| `json` | Output is formatted as JSON |
| `jsonPretty` | Outputs a human-readable JSON with indented by 2 spaces |
| `upper` | Convert string to uppercase |
| `lower` | Convert string to lowercase |
| `split` | Splits strings given by `string` and `sep` |
| `join` | Joins strings by given separator |
| `title` | Convert the first letter to uppercase of a string |

# Example

In this example we are going to printout details of few persons using teks.

```go
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
		{"John Doe", 34, "Somewhere on earth"},
		{"Spongebob", 30, "Under Sea"},
		{"Harry Potter", 30, "4 Privet Drive, Little Whinging, Surrey"},
	}

	// create new context
	ctx := teks.NewContext(os.Stdout, format)

	// create a renderer function to match signature defined in teks.Renderer
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

```

Now run program as follows

```
➜ go run simple.go 
Name                Age
Kasun               24
John Doe            34
Spongebob           30
Harry Potter        30
```

Let's pretty print Name and Address in tabular format
```
➜ go run simple.go --format "table {{.Name}}\t{{.Address}}"
Name                Address
Kasun               Earth
John Doe            Somewhere on earth
Spongebob           Under Sea
Harry Potter        4 Privet Drive, Little Whinging, Surrey
```

Let's make Name UPPERCASE
```
➜ go run simple.go --format "table {{upper .Name}}\t{{.Address}}"
NAME                Address
KASUN               Earth
JOHN DOE            Somewhere on earth
SPONGEBOB           Under Sea
HARRY POTTER        4 Privet Drive, Little Whinging, Surrey
```

You can change behavior of these headers by providing custom HeaderFuncs.

[![asciicast](https://asciinema.org/a/249879.svg)](https://asciinema.org/a/249879)

# Contribution

All contributions are welcome. Raise an Issue or a Pull Request
