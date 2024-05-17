package client

import (
	"crypto/tls"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// Response represents a Gemini Response
type Response struct {
	// Data contains the raw data returned from the server.
	// TODO: implement function to pipe this straight to Gemtext parser
	Data []byte
	// MetaData contains the metadata of the response.
	// If StatusCode == 20, it is the MIME type associated with the data
	MetaData string
	// StatusCode contains the status code returned by the server.
	// If this is not 20, Data should be considered to be empty
	StatusCode int
}

// Request sends a Gemini request to address
// tlsConfig will be removed in a future update. Per the tls.Client documentation,
//
// The config cannot be nil: users must set either ServerName or
// InsecureSkipVerify in the config.
//
// TODO: Does not currently handle redirects.
func Request(address string, tlsConfig *tls.Config) (*Response, error) {
	parsedURL, err := url.Parse(address)
	if err != nil {
		return nil, err
	}
	connInsecure, err := net.Dial("tcp", parsedURL.Host+":1965")
	if err != nil {
		return nil, err
	}

	conn := tls.Client(connInsecure, tlsConfig)

	defer func(conn *tls.Conn) {
		err := conn.Close()
		if err != nil {
			log.Errorf("%v", err)
		}
	}(conn)

	_, err = conn.Write([]byte(address + "\r\n"))
	if err != nil {
		return nil, err
	}

	responseData, err := io.ReadAll(conn)
	if err != nil {
		return nil, err
	}

	header := strings.Split(string(responseData), "\r\n")[0]
	headerParts := strings.Split(header, " ")
	statusCode, err := strconv.Atoi(headerParts[0])
	if err != nil {
		return nil, err
	}
	metaData := strings.Join(headerParts[1:], " ")

	return &Response{
		StatusCode: statusCode,
		Data:       responseData[len(header)+2:],
		MetaData:   metaData,
	}, nil
}
