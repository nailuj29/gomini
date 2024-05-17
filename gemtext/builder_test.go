package gemtext

import "testing"

func TestBuilder_AddHeader1Line(t *testing.T) {
	b := NewBuilder()

	b.AddHeader1Line("Hello")

	if b.Get() != "# Hello" {
		t.Errorf("got %q, want %q", b.Get(), "# Hello")
	}
}

func TestBuilder_AddHeader2Line(t *testing.T) {
	b := NewBuilder()

	b.AddHeader2Line("Hello")

	if b.Get() != "## Hello" {
		t.Errorf("got %q, want %q", b.Get(), "## Hello")
	}
}

func TestBuilder_AddHeader3Line(t *testing.T) {
	b := NewBuilder()

	b.AddHeader3Line("Hello")

	if b.Get() != "### Hello" {
		t.Errorf("got %q, want %q", b.Get(), "### Hello")
	}
}

func TestBuilder_AddTextLine(t *testing.T) {
	b := NewBuilder()

	b.AddTextLine("Hello")

	if b.Get() != "Hello" {
		t.Errorf("got %q, want %q", b.Get(), "Hello")
	}
}

func TestBuilder_AddLinkLine(t *testing.T) {
	b := NewBuilder()

	b.AddLinkLine("gemini://geminiprotocol.net")
	if b.Get() != "=> gemini://geminiprotocol.net" {
		t.Errorf("got %q, want %q", b.Get(), "=> gemini://geminiprotocol.net")
	}

	b = NewBuilder()

	b.AddLinkLine("gemini://geminiprotocol.net", "Gemini Protocol Website")
	if b.Get() != "=> gemini://geminiprotocol.net Gemini Protocol Website" {
		t.Errorf("got %q, want %q", b.Get(), "=> gemini://geminiprotocol.net Gemini Protocol Website")
	}
}

func TestBuilder_AddPreformattedText(t *testing.T) {
	b := NewBuilder()

	b.AddPreformattedText("Hello\r\nWorld")
	if b.Get() != "```\r\nHello\r\nWorld\r\n```" {
		t.Errorf("got %q, want %q", b.Get(), "```\nHello\nWorld\n```")
	}
}

func TestBuilder_AddUnorderedList(t *testing.T) {
	b := NewBuilder()

	b.AddUnorderedList([]string{
		"item 1",
		"item 2",
		"item 3",
	})

	if b.Get() != "* item 1\r\n* item 2\r\n* item 3" {
		t.Errorf("got %q, want %q", b.Get(), "* item 1\r\n* item 2\r\n* item 3")
	}
}

func TestBuilder_AddQuoteLine(t *testing.T) {
	b := NewBuilder()

	b.AddQuoteLine("Hello")

	if b.Get() != "> Hello" {
		t.Errorf("got %q, want %q", b.Get(), "> Hello")
	}
}
