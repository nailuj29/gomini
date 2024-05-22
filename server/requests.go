package server

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/url"
	"os"
)

// Request wraps a Gemini request.
type Request struct {
	// URI contains a [url.URL] object corresponding to the URL of the request.
	URI url.URL
	// Params contains a map of URL params passed into the request. Nil if there are no params.
	Params     map[string]string
	conn       *tls.Conn
	terminated bool
}

// TitanRequest wraps a Titan request.
//
// Contains [Request], so can be used like a Gemini request as well.
type TitanRequest struct {
	Request
	// Token contains the answer to the security question, or an empty string
	Token string
	// MIMEType contains the MIME type of the data. Defaults to "text/gemini"
	MIMEType string
	// Body contains the data sent by the client
	Body []byte
}

// Gemtext responds using a gemtext string and status code 20.
// After calling this method, the [Request] has been terminated.
func (r *Request) Gemtext(source string) error {
	if r.terminated {
		return errors.New("already responded")
	}

	_, err := r.conn.Write([]byte("20 text/gemini\r\n" + source))
	if err != nil {
		return err
	}

	r.terminated = true

	return nil
}

// GemtextFile responds using gemtext from a file and status code 20.
// After calling this method, the [Request] has been terminated.
func (r *Request) GemtextFile(path string) error {
	if r.terminated {
		return errors.New("already responded")
	}

	source, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return r.Gemtext(string(source))
}

// Error responds with an error code and message
// After calling this method, the [Request] has been terminated.
func (r *Request) Error(code int, message string) error {
	_, err := r.conn.Write([]byte(fmt.Sprintf("%d %s\r\n", code, message)))

	return err
}

// GetClientCertificates retrieves the client certificate(s) for the [Request]
func (r *Request) GetClientCertificates() []*x509.Certificate {
	return r.conn.ConnectionState().PeerCertificates
}

// TODO: Cert signatures
