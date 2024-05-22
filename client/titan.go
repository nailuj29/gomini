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

// TitanRequest sends a request to a Titan server
//
// If token or mime is not desired, an empty string can be passed.
// The same caveat for tlsConfig as [Request] applies.
func TitanRequest(address string, tlsConfig *tls.Config, body []byte, token string, mime string) (*Response, error) {
	if mime == "" {
		mime = "text/gemini"
	}

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

	uri := address + ";size=" + strconv.Itoa(len(body)) + ";mime=" + mime
	if token != "" {
		uri += ";token=" + token
	}

	_, err = conn.Write(append([]byte(uri+"\r\n"), body...))
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
