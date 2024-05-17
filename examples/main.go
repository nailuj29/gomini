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
		err := request.GemtextFile("index.gmi")
		if err != nil {
			log.Fatalf("could not respond: %v", err)
		}
	})

	s.RegisterHandler("/test1", func(request server.Request) {
		err := request.Gemtext("# Test 1!\r\nThis is the first test page")
		if err != nil {
			log.Fatalf("could not respond: %v", err)
		}
	})

	s.RegisterHandler("/test2", func(request server.Request) {
		err := request.Gemtext("# Test 2!\r\nThis is the second test page")
		if err != nil {
			log.Fatalf("could not respond: %v", err)
		}
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

		err := request.Gemtext(b.Get())
		if err != nil {
			log.Fatalf("could not respond: %v", err)
		}
	})

	s.RegisterHandler("/secure", func(request server.Request) {
		certs := request.GetClientCertificates()

		if len(certs) == 0 {
			err := request.Error(60, "Cert required")
			if err != nil {
				log.Fatalf("could not respond: %v", err)
			}
			return
		}

		err := request.Gemtext("# Secure page\r\nWelcome!")
		if err != nil {
			log.Fatalf("could not respond: %v", err)
		}
	})

	s.RegisterHandler("/dynamic/:dynamic", func(request server.Request) {
		param := request.Params["dynamic"]

		b := gemtext.NewBuilder()
		b.
			AddHeader1Line("# Dynamic Route").
			AddTextLine(param)

		err := request.Gemtext(b.Get())
		if err != nil {
			log.Fatalf("could not respond: %v", err)
		}
	})

	err = s.ListenAndServe("localhost", &config)
	if err != nil {
		log.Fatalf("could not respond: %v", err)
	}
}
