package main

import (
	"crypto/tls"
	"log"

	"github.com/nailuj29gaming/gomini/server"
)

func main() {
	cer, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cer}}

	s := server.New()

	s.AddHandler("/", func(request server.Request) {
		request.GemtextFile("index.gmi")
	})

	s.AddHandler("/test1", func(request server.Request) {
		request.Gemtext("# Test 1!\r\nThis is the first test page")
	})

	s.AddHandler("/test2", func(request server.Request) {
		request.Gemtext("# Test 2!\r\nThis is the second test page")
	})

	s.ListenAndServe("localhost", &config)
}
