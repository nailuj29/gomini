package gemtext

import (
	"fmt"
	"strings"
)

// A Builder is an object capable of building a Gemtext document
type Builder struct {
	lines []string
}

// NewBuilder creates a new [Builder]
func NewBuilder() Builder {
	return Builder{
		lines: make([]string, 0),
	}
}

// AddTextLine adds a new text line (i.e. no additional formatting) to the [Builder]
func (b *Builder) AddTextLine(line string) *Builder {
	b.lines = append(b.lines, line)
	return b
}

// AddLinkLine adds a new link line.
//
// While more than two args can be passed to it, any additional arguments will be ignored for the time being
func (b *Builder) AddLinkLine(url string, name ...string) *Builder {
	// TODO: Vararg is a huge bandaid, fix it :D
	if len(name) == 0 {
		b.AddTextLine(fmt.Sprintf("=> %s", url))
	} else {
		b.AddTextLine(fmt.Sprintf("=> %s %s", url, name[0]))
	}

	return b
}

// AddPreformattedText adds a block of preformatted text
func (b *Builder) AddPreformattedText(text string) *Builder {
	b.AddTextLine("```")
	b.AddTextLine(text)
	b.AddTextLine("```")

	return b
}

// AddHeader1Line adds a top level header line
func (b *Builder) AddHeader1Line(text string) *Builder {
	b.AddTextLine(fmt.Sprintf("# %s", text))

	return b
}

// AddHeader2Line adds a secondary level header line
func (b *Builder) AddHeader2Line(text string) *Builder {
	b.AddTextLine(fmt.Sprintf("## %s", text))

	return b
}

// AddHeader3Line adds a tertiary level header line
func (b *Builder) AddHeader3Line(text string) *Builder {
	b.AddTextLine(fmt.Sprintf("### %s", text))

	return b
}

// AddUnorderedList adds an unordered list
func (b *Builder) AddUnorderedList(items []string) *Builder {
	for _, item := range items {
		b.AddTextLine(fmt.Sprintf("* %s", item))
	}

	return b
}

// AddQuoteLine adds a blockquoted line
func (b *Builder) AddQuoteLine(text string) *Builder {
	b.AddTextLine(fmt.Sprintf("> %s", text))

	return b
}

// Get gets the Gemtext content of the Builder
func (b *Builder) Get() string {
	return strings.Join(b.lines, "\r\n")
}
