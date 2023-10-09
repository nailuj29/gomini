package gemtext

import (
	"fmt"
	"strings"
)

type Builder struct {
	lines []string
}

func NewBuilder() Builder {
	return Builder{
		lines: make([]string, 0),
	}
}

func (b *Builder) AddTextLine(line string) *Builder {
	b.lines = append(b.lines, line)
	return b
}

func (b *Builder) AddLinkLine(url string, name string) *Builder {
	if len(name) == 0 {
		b.AddTextLine(fmt.Sprintf("=> %s", url))
	} else {
		b.AddTextLine(fmt.Sprintf("=> %s", name[0]))
	}

	return b
}

func (b *Builder) AddPreformattedText(text string) *Builder {
	b.AddTextLine("```")
	b.AddTextLine(text)
	b.AddTextLine("```")

	return b
}

func (b *Builder) AddHeader1Line(text string) *Builder {
	b.AddTextLine(fmt.Sprintf("# %s", text))

	return b
}

func (b *Builder) AddHeader2Line(text string) *Builder {
	b.AddTextLine(fmt.Sprintf("## %s", text))

	return b
}

func (b *Builder) AddHeader3Line(text string) *Builder {
	b.AddTextLine(fmt.Sprintf("### %s", text))

	return b
}

func (b *Builder) AddUnorderedList(items []string) *Builder {
	for _, item := range items {
		b.AddTextLine(fmt.Sprintf("* %s", item))
	}

	return b
}

func (b *Builder) AddQuoteLine(text string) *Builder {
	b.AddTextLine(fmt.Sprintf("> %s", text))

	return b
}

func (b *Builder) Get() string {
	return strings.Join(b.lines, "\r\n")
}
