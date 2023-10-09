package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/url"
	"os"
)

type Request struct {
	URI  url.URL
	conn *tls.Conn
}

func (r *Request) Gemtext(source string) error {
	_, err := r.conn.Write([]byte("20 text/gemini\r\n" + source))

	return err
}

func (r *Request) GemtextFile(path string) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return r.Gemtext(string(source))
}

func (r *Request) Error(code int, message string) error {
	_, err := r.conn.Write([]byte(fmt.Sprintf("%d %s", code, message)))

	return err
}

func (r *Request) GetClientCertificates() []*x509.Certificate {
	return r.conn.ConnectionState().PeerCertificates
}
