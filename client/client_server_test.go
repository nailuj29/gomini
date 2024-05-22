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
