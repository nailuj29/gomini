package gemtext

import (
	"errors"
	"strings"
)

// LineType is an enumeration representing a Gemtext line type
type LineType int

// Line is any object representing a Gemtext line
type Line interface {
	// Type returns the type of the line
	Type() LineType
}

// TextLine is a plain text line, with no additional formatting
type TextLine struct {
	// Text is the line's text
	Text string
}

// LinkLine is a hyperlink line
type LinkLine struct {
	// Destination is the line's destination
	Destination string // TODO: Maybe make this a url.URL?
	// Text is the line's (optional) body text
	Text string
}

// PreformattedText is a block of preformatted text
type PreformattedText struct {
	// AltText is the block's alt text; often used for language syntax highlighting
	AltText string
	// Body is the body of the block
	Body string
}

// Header1Line is a primary header
type Header1Line struct {
	// Text is the text of the header
	Text string
}

// Header2Line is a secondary header
type Header2Line struct {
	// Text is the text of the header
	Text string
}

// Header3Line is a tertiary header
type Header3Line struct {
	// Text is the text of the header
	Text string
}

// ListItemLine is a item in an unordered list
type ListItemLine struct {
	// Text is the text of the list item
	Text string
}

// QuoteLine is a blockquoted line
type QuoteLine struct {
	// Text is the text within the quote
	Text string
}

const (
	// Text marks a plain text line, with no additional formatting
	Text LineType = iota
	// Link marks a hyperlink line
	Link
	// Preformatted marks a block of preformatted text
	Preformatted
	// Header1 marks a primary header
	Header1
	// Header2 marks a secondary header
	Header2
	// Header3 marks a tertiary header
	Header3
	// ListItem marks an item in an unordered list
	ListItem
	// Quote marks a blockquoted line
	Quote
)

func (t TextLine) Type() LineType {
	return Text
}

func (l LinkLine) Type() LineType {
	return Link
}

func (p PreformattedText) Type() LineType {
	return Preformatted
}

func (h Header1Line) Type() LineType {
	return Header1
}

func (h Header2Line) Type() LineType {
	return Header2
}

func (h Header3Line) Type() LineType {
	return Header3
}

func (l ListItemLine) Type() LineType {
	return ListItem
}

func (q QuoteLine) Type() LineType {
	return Quote
}

func Parse(source string) ([]Line, error) {
	sourceLines := strings.Split(source, "\r\n")
	lines := make([]Line, 0)
	preformattingToggled := false
	preAltText := ""
	preLines := make([]string, 0)

	for _, l := range sourceLines {
		if preformattingToggled && !strings.HasPrefix(l, "```") {
			preLines = append(preLines, l)
			continue
		}
		if strings.HasPrefix(l, "=>") {
			ll := strings.Fields(l)
			if len(ll) < 2 {
				return nil, errors.New("missing destination for link line")
			}

			dest := ll[1]
			if len(ll) == 2 {
				lines = append(lines, LinkLine{
					Destination: dest,
					Text:        "",
				})
			} else {
				whiteSpaceDestAndText := strings.TrimPrefix(l, "=>")
				destAndText := strings.TrimLeft(whiteSpaceDestAndText, " \t")
				text := strings.TrimPrefix(destAndText, dest)
				lines = append(lines, LinkLine{
					Destination: dest,
					Text:        strings.TrimLeft(text, " \t"),
				})
			}
		} else if strings.HasPrefix(l, "```") {
			if preformattingToggled {
				preformattingToggled = false
				lines = append(lines, PreformattedText{
					Body:    strings.Join(preLines, "\r\n"),
					AltText: preAltText,
				})

				continue
			}

			preAltText = strings.TrimPrefix(l, "```")
			preformattingToggled = true
		} else if strings.HasPrefix(l, "###") { // Must work backwards, or everything is a Header1
			text := strings.TrimPrefix(l, "###")
			lines = append(lines, Header3Line{
				Text: strings.TrimLeft(text, " \t"),
			})
		} else if strings.HasPrefix(l, "##") {
			text := strings.TrimPrefix(l, "##")
			lines = append(lines, Header2Line{
				Text: strings.TrimLeft(text, " \t"),
			})
		} else if strings.HasPrefix(l, "#") {
			text := strings.TrimPrefix(l, "#")
			lines = append(lines, Header1Line{
				Text: strings.TrimLeft(text, " \t"),
			})
		} else if strings.HasPrefix(l, "*") {
			text := strings.TrimPrefix(l, "*")
			lines = append(lines, ListItemLine{
				Text: strings.TrimLeft(text, " \t"),
			})
		} else if strings.HasPrefix(l, ">") {
			text := strings.TrimPrefix(l, ">")
			lines = append(lines, QuoteLine{
				Text: strings.TrimLeft(text, " \t"),
			})
		} else {
			lines = append(lines, TextLine{
				Text: l,
			})
		}
	}

	if preformattingToggled {
		return nil, errors.New("unclosed preformatting block")
	}

	return lines, nil
}
