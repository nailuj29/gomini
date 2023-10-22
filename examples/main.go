package main

import (
	"crypto/tls"
	"github.com/nailuj29/gomini/gemtext"
	"github.com/nailuj29/gomini/server"
	"log"
)

func main() {
	cer, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cer}, ClientAuth: tls.RequestClientCert}

	s := server.New()

	s.RegisterHandler("/", func(request server.Request) {
		request.GemtextFile("index.gmi")
	})

	s.RegisterHandler("/test1", func(request server.Request) {
		request.Gemtext("# Test 1!\r\nThis is the first test page")
	})

	s.RegisterHandler("/test2", func(request server.Request) {
		request.Gemtext("# Test 2!\r\nThis is the second test page")
	})

	s.RegisterHandler("/gemtext", func(request server.Request) {
		b := gemtext.NewBuilder()
		b.
			AddHeader1Line("Gemtext").
			AddHeader2Line("Level 2").
			AddHeader3Line("Builder").
			AddTextLine("Text Lines").
			AddPreformattedText("Oh cool, code!").
			AddLinkLine("gemini://localhost/", "Go Home").
			AddTextLine("").
			AddQuoteLine("Please stop making up stuff I said").
			AddTextLine("- Sun Tzu, Art of War").
			AddUnorderedList([]string{
				"Item 1",
				"Item 2",
				"Item 3",
			})

		request.Gemtext(b.Get())
	})

	s.RegisterHandler("/secure", func(request server.Request) {
		certs := request.GetClientCertificates()

		if len(certs) == 0 {
			request.Error(60, "Cert required")
			return
		}

		request.Gemtext("# Secure page\r\nWelcome!")
	})

	s.RegisterHandler("/dynamic/:dynamic", func(request server.Request) {
		param := request.Params["dynamic"]

		b := gemtext.NewBuilder()
		b.
			AddHeader1Line("# Dynamic Route").
			AddTextLine(param)

		request.Gemtext(b.Get())
	})

	s.ListenAndServe("localhost", &config)
}
