package server

import (
	"crypto/tls"
	"crypto/x509"
	"net/url"
	"os"
)

type Request struct {
	URI  url.URL
	conn *tls.Conn
}

func (r *Request) Gemtext(source string) error {
	_, err := r.conn.Write([]byte("20 text/gemini\r\n" + source))
	if err != nil {
		return err
	}

	return nil
}

func (r *Request) GemtextFile(path string) error {
	source, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return r.Gemtext(string(source))
}

func (r *Request) GetClientCertificates() []*x509.Certificate {
	return r.conn.ConnectionState().PeerCertificates
}
