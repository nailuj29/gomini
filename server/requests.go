package server

import (
	"crypto/tls"
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

	_, err = r.conn.Write([]byte("20 text/gemini\r\n" + string(source)))
	if err != nil {
		return err
	}

	return nil
}
