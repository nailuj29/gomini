package client

import (
	"crypto/tls"
	"github.com/nailuj29/gomini/gemtext"
	"github.com/nailuj29/gomini/server"
	"testing"
)

func TestBasicRequest(t *testing.T) {
	cer, err := tls.LoadX509KeyPair("../examples/cert.pem", "../examples/key.pem")
	if err != nil {
		t.Fatal(err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cer}, ClientAuth: tls.RequestClientCert}

	s := server.New()

	s.RegisterHandler("/", func(r server.Request) {
		b := gemtext.NewBuilder()
		b.AddHeader1Line("Hello, World!")

		r.Gemtext(b.Get())
	})

	go func() {
		err := s.ListenAndServe("localhost", &config)
		if err != nil {
			t.Fatal(err)
		}
	}()

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
