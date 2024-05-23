package client

import (
	"crypto/tls"
	"github.com/nailuj29/gomini/gemtext"
	"github.com/nailuj29/gomini/server"
	"testing"
)

func TestBasicRequestResponse(t *testing.T) {
	cer, err := tls.LoadX509KeyPair("../examples/cert.pem", "../examples/key.pem")
	if err != nil {
		t.Fatal(err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cer}, ClientAuth: tls.RequestClientCert}

	s := server.New()

	s.RegisterHandler("/", func(r server.Request) {
		b := gemtext.NewBuilder()
		b.AddHeader1Line("Hello, World!")

		err := r.Gemtext(b.Get())
		if err != nil {
			t.Errorf("handler failed to respond to request: %v", err)
		}
	})

	go func() {
		s.ListenAndServe("localhost", &config)
	}()

	defer func(s *server.Server) {
		err := s.Close()
		if err != nil {
			t.Errorf("Could not close server: %v", err)
		}
	}(s)

	clientConfig := tls.Config{InsecureSkipVerify: true}
	response, err := Request("gemini://localhost/", &clientConfig)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != 20 {
		t.Fatalf("Response status code is %d", response.StatusCode)
	}

	if response.MetaData != "text/gemini" {
		t.Fatalf("Response meta data is %s", response.MetaData)
	}

	parsed, err := gemtext.Parse(string(response.Data))
	if err != nil {
		t.Fatal(err)
	}

	headerLine, ok := parsed[0].(gemtext.Header1Line)
	if !ok {
		t.Fatalf("First line not a valid header line. Response data = %s", string(response.Data))
	}

	if headerLine.Text != "Hello, World!" {
		t.Fatalf("First line text is %s", headerLine.Text)
	}
}

func TestDynamicPathRequestResponse(t *testing.T) {
	cer, err := tls.LoadX509KeyPair("../examples/cert.pem", "../examples/key.pem")
	if err != nil {
		t.Fatal(err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cer}, ClientAuth: tls.RequestClientCert}

	s := server.New()

	s.RegisterHandler("/:param", func(r server.Request) {
		err := r.Gemtext(r.Params["param"])
		if err != nil {
			t.Errorf("handler failed to respond to request: %v", err)
		}
	})

	go func() {
		s.ListenAndServe("localhost", &config)
	}()

	defer func(s *server.Server) {
		err := s.Close()
		if err != nil {
			t.Errorf("Could not close server: %v", err)
		}
	}(s)

	clientConfig := tls.Config{InsecureSkipVerify: true}
	response, err := Request("gemini://localhost/foo-bar", &clientConfig)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != 20 {
		t.Fatalf("Response status code is %d", response.StatusCode)
	}

	if response.MetaData != "text/gemini" {
		t.Fatalf("Response meta data is %s", response.MetaData)
	}

	parsed, err := gemtext.Parse(string(response.Data))
	if err != nil {
		t.Fatal(err)
	}

	textLine, ok := parsed[0].(gemtext.TextLine)
	if !ok {
		t.Fatalf("First line not a text line. Response data = %s", string(response.Data))
	}

	if textLine.Text != "foo-bar" {
		t.Fatalf("First line text is %s", textLine.Text)
	}
}

func TestTitanRequestResponse(t *testing.T) {
	cer, err := tls.LoadX509KeyPair("../examples/cert.pem", "../examples/key.pem")
	if err != nil {
		t.Fatal(err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cer}, ClientAuth: tls.RequestClientCert}

	s := server.New()

	s.RegisterTitanHandler("/", func(r server.TitanRequest) {
		b := gemtext.NewBuilder()
		b.AddPreformattedText(string(r.Body))
		b.AddHeader1Line(r.Token)
		b.AddHeader2Line(r.MIMEType)

		err := r.Gemtext(b.Get())
		if err != nil {
			t.Errorf("handler failed to respond to request: %v", err)
		}
	})

	go func() {
		s.ListenAndServe("localhost", &config)
	}()

	defer func(s *server.Server) {
		err := s.Close()
		if err != nil {
			t.Errorf("Could not close server: %v", err)
		}
	}(s)

	clientConfig := tls.Config{InsecureSkipVerify: true}
	response, err := TitanRequest("titan://localhost/", &clientConfig, []byte("Hello"), "tokenTester", "text/gemini")
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != 20 {
		t.Fatalf("Response status code is %d", response.StatusCode)
	}

	if response.MetaData != "text/gemini" {
		t.Fatalf("Response meta data is %s", response.MetaData)
	}

	parsed, err := gemtext.Parse(string(response.Data))
	if err != nil {
		t.Fatal(err)
	}

	lines := []gemtext.Line{
		gemtext.PreformattedText{
			Body: "Hello",
		},
		gemtext.Header1Line{
			Text: "tokenTester",
		},
		gemtext.Header2Line{
			Text: "text/gemini",
		},
	}

	for i, line := range lines {
		parsedLine := parsed[i]

		if parsedLine.Type() != line.Type() {
			t.Fatalf("Line type is %d, want %d", parsedLine.Type(), line.Type())
		}

		switch parsedLine.Type() {
		case gemtext.Header1:
			if parsedLine.(gemtext.Header1Line).Text != line.(gemtext.Header1Line).Text {
				t.Fatalf("Expected %s, got %s", line.(gemtext.Header1Line).Text, parsedLine.(gemtext.Header1Line).Text)
			}
		case gemtext.Header2:
			if parsedLine.(gemtext.Header2Line).Text != line.(gemtext.Header2Line).Text {
				t.Fatalf("Expected %s, got %s", line.(gemtext.Header2Line).Text, parsedLine.(gemtext.Header2Line).Text)
			}
		case gemtext.Preformatted:
			if parsedLine.(gemtext.PreformattedText).Body != line.(gemtext.PreformattedText).Body {
				t.Fatalf("Expected %s, got %s", line.(gemtext.PreformattedText).Body, parsedLine.(gemtext.PreformattedText).Body)
			}
		default:
			t.Fatalf("Unexpected Line type %d", parsedLine.Type())
		}
	}
}

func TestInputRequestResponse(t *testing.T) {
	cer, err := tls.LoadX509KeyPair("../examples/cert.pem", "../examples/key.pem")
	if err != nil {
		t.Fatal(err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cer}, ClientAuth: tls.RequestClientCert}

	s := server.New()

	s.RegisterHandler("/", func(r server.Request) {
		inp, err := r.RequestInput("Input please: ")
		if err != nil {
			t.Fatal(err)
		}

		if inp == "" {
			return
		}

		r.Gemtext(inp)
	})

	go func() {
		s.ListenAndServe("localhost", &config)
	}()

	defer func(s *server.Server) {
		err := s.Close()
		if err != nil {
			t.Errorf("Could not close server: %v", err)
		}
	}(s)

	clientConfig := tls.Config{InsecureSkipVerify: true}
	response, err := Request("gemini://localhost/", &clientConfig)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != 10 {
		t.Fatalf("Response status code is %d", response.StatusCode)
	}

	if response.MetaData != "Input please: " {
		t.Fatalf("Response meta data is %s", response.MetaData)
	}

	response, err = Request("gemini://localhost/?Hello%2C%20World%21", &clientConfig)
	if err != nil {
		t.Fatal(err)
	}

	if response.StatusCode != 20 {
		t.Fatalf("Response status code is %d", response.StatusCode)
	}

	if response.MetaData != "text/gemini" {
		t.Fatalf("Response meta data is %s", response.MetaData)
	}

	parsed, err := gemtext.Parse(string(response.Data))
	if err != nil {
		t.Fatal(err)
	}
	textLine, ok := parsed[0].(gemtext.TextLine)
	if !ok {
		t.Fatalf("First line not a text line. Response data = %s", string(response.Data))
	}

	if textLine.Text != "Hello, World!" {
		t.Fatalf("First line text is %s", textLine.Text)
	}
}
