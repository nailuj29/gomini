package gemtext_test

import (
	"github.com/nailuj29/gomini/gemtext"
	"testing"
)

func TestParsingText(t *testing.T) {
	source := "Text Line 1\r\nText Line 2"
	parsed, err := gemtext.Parse(source)
	if err != nil {
		t.Fatalf("Error in parsing: %s", err.Error())
	}

	line1 := parsed[0]
	if line1.Type() != gemtext.Text {
		t.Fatalf("Didn't get text line for line 1")
	}

	textLine1 := line1.(gemtext.TextLine)
	if textLine1.Text != "Text Line 1" {
		t.Errorf("Wrong text in line 1, expected \"%s\", got \"%s\"", "Text Line 1", textLine1.Text)
	}

	line2 := parsed[1]
	if line2.Type() != gemtext.Text {
		t.Fatalf("Didn't get text line for line 2")
	}

	textLine2 := line2.(gemtext.TextLine)
	if textLine2.Text != "Text Line 2" {
		t.Errorf("Wrong text in line 2, expected \"%s\", got \"%s\"", "Text Line 2", textLine2.Text)
	}
}

func TestParsingLinks(t *testing.T) {
	source := "=> gemini://example.com\r\n=> gemini://example.com Example link\r\n=>    gemini://example.com \tWeird Spacing"

	parsed, err := gemtext.Parse(source)
	if err != nil {
		t.Fatalf("Error in parsing: %s", err.Error())
	}

	line1 := parsed[0]
	if line1.Type() != gemtext.Link {
		t.Fatalf("Didn't get link line for line 1")
	}

	linkLine1 := line1.(gemtext.LinkLine)
	if linkLine1.Destination != "gemini://example.com" {
		t.Errorf("Got wrong link on line 1: \"%s\"", linkLine1.Destination)
	}

	if linkLine1.Text != "" {
		t.Errorf("Got text on line 1: \"%s\"", linkLine1.Text)
	}

	line2 := parsed[1]
	if line2.Type() != gemtext.Link {
		t.Fatalf("Didn't get link line for line 2")
	}

	linkLine2 := line2.(gemtext.LinkLine)
	if linkLine2.Destination != "gemini://example.com" {
		t.Errorf("Got wrong link on line 2: \"%s\"", linkLine2.Destination)
	}

	if linkLine2.Text != "Example link" {
		t.Errorf("Got wrong text on line 2. Expected \"%s\", got \"%s\"", "Example link", linkLine2.Text)
	}

	line3 := parsed[2]
	if line3.Type() != gemtext.Link {
		t.Fatalf("Didn't get link line for line 3")
	}

	linkLine3 := line3.(gemtext.LinkLine)
	if linkLine3.Destination != "gemini://example.com" {
		t.Errorf("Got wrong link on line 3: \"%s\"", linkLine3.Destination)
	}

	if linkLine3.Text != "Weird Spacing" {
		t.Errorf("Got wrong text on line 3. Expected \"%s\", got \"%s\"", "Weird Spacing", linkLine3.Text)
	}
}

func TestParsingPreBlocks(t *testing.T) {
	source := "```\r\ntext\r\nline2\r\n```\r\n```alt-text\r\nbody\r\n```"

	parsed, err := gemtext.Parse(source)
	if err != nil {
		t.Fatalf("Error in parsing: %s", err.Error())
	}

	line1 := parsed[0]
	if line1.Type() != gemtext.Preformatted {
		t.Fatalf("Didn't get preformatted block for line 1")
	}

	preBlock1 := line1.(gemtext.PreformattedText)
	if preBlock1.AltText != "" {
		t.Errorf("Got alt text for line 1: \"%s\"", preBlock1.AltText)
	}

	if preBlock1.Body != "text\r\nline2" {
		t.Errorf("Got wrong body text for line 1. Expected: \"%s, got: \"%s\"", "text\\r\\nline2", preBlock1.Body)
	}

	line2 := parsed[1]
	if line2.Type() != gemtext.Preformatted {
		t.Fatalf("Didn't get preformatted block for line 2")
	}

	preBlock2 := line2.(gemtext.PreformattedText)
	if preBlock2.AltText != "alt-text" {
		t.Errorf("Got wrong alt text for line 2. Expected: \"%s\", got \"%s\"", "alt-text", preBlock1.AltText)
	}

	if preBlock1.Body != "text\r\nline2" {
		t.Errorf("Got wrong body text for line 1. Expected: \"%s\", got: \"%s\"", "text\\r\\nline2", preBlock1.Body)
	}
}

func TestParsingHeaders(t *testing.T) {
	source := "# Header 1\r\n## Header 2\r\n### Header 3"

	parsed, err := gemtext.Parse(source)
	if err != nil {
		t.Fatalf("Error in parsing: %s", err.Error())
	}

	line1 := parsed[0]
	if line1.Type() != gemtext.Header1 {
		t.Errorf("Didn't get header level 1 for line 1")
	}

	header1 := line1.(gemtext.Header1Line)
	if header1.Text != "Header 1" {
		t.Errorf("Got wrong text for line 1. Expected: \"%s\", for \"%s\"", "Header 1", header1.Text)
	}

	line2 := parsed[1]
	if line2.Type() != gemtext.Header2 {
		t.Errorf("Didn't get header level 2 for line 2")
	}

	header2 := line2.(gemtext.Header2Line)
	if header2.Text != "Header 2" {
		t.Errorf("Got wrong text for line 2. Expected: \"%s\", for \"%s\"", "Header 2", header2.Text)
	}

	line3 := parsed[2]
	if line3.Type() != gemtext.Header3 {
		t.Errorf("Didn't get header level 3 for line 3")
	}

	header3 := line3.(gemtext.Header3Line)
	if header3.Text != "Header 3" {
		t.Errorf("Got wrong text for line 3. Expected: \"%s\", for \"%s\"", "Header 3", header3.Text)
	}
}

func TestParsingLists(t *testing.T) {
	source := "* Item 1\r\n* Item 2"

	parsed, err := gemtext.Parse(source)
	if err != nil {
		t.Fatalf("Error in parsing: %s", err.Error())
	}

	line1 := parsed[0]
	if line1.Type() != gemtext.ListItem {
		t.Errorf("Didn't get list item for line 1")
	}

	item1 := line1.(gemtext.ListItemLine)
	if item1.Text != "Item 1" {
		t.Errorf("Got wrong text for line 1. Expected: \"%s\", for \"%s\"", "Header 1", item1.Text)
	}

	line2 := parsed[1]
	if line2.Type() != gemtext.ListItem {
		t.Errorf("Didn't get list item for line 2")
	}

	item2 := line2.(gemtext.ListItemLine)
	if item2.Text != "Item 2" {
		t.Errorf("Got wrong text for line 2. Expected: \"%s\", for \"%s\"", "Header 2", item2.Text)
	}
}

func TestParsingQuotes(t *testing.T) {
	source := "> Quotation 1\r\n> Quotation 2"

	parsed, err := gemtext.Parse(source)
	if err != nil {
		t.Fatalf("Error in parsing: %s", err.Error())
	}

	line1 := parsed[0]
	if line1.Type() != gemtext.Quote {
		t.Errorf("Didn't get quote for line 1")
	}

	quote1 := line1.(gemtext.QuoteLine)
	if quote1.Text != "Quotation 1" {
		t.Errorf("Got wrong text for line 1. Expected: \"%s\", for \"%s\"", "Quotation 1", quote1.Text)
	}

	line2 := parsed[1]
	if line2.Type() != gemtext.Quote {
		t.Errorf("Didn't get quote for line 2")
	}

	quote2 := line2.(gemtext.QuoteLine)
	if quote2.Text != "Quotation 2" {
		t.Errorf("Got wrong text for line 2. Expected: \"%s\", for \"%s\"", "Quotation 2", quote2.Text)
	}
}
