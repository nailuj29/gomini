// Package server provides utilities for creating a functioning and useful Gemini server
package server

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

// A Handler is a function to handle a Request by calling its various methods.
// The function is called when a request that it can handle, as outlined in Server.RegisterHandler
type Handler func(request Request)

// A Server contains information required to run a TCP/TLS service capable of serving Gemini content over the internet
type Server struct {
	staticRoutes  map[string]Handler
	dynamicRoutes []route
	listener      net.Listener
	addr          string
}

type route struct {
	regex   *regexp.Regexp
	handler Handler
}

// New creates a new Server
func New() *Server {
	return &Server{}
}

// RegisterHandler sets up a Handler to handle any Request that comes to a path
func (s *Server) RegisterHandler(path string, handler Handler) {
	if !strings.ContainsRune(path, ':') {
		if s.staticRoutes == nil {
			s.staticRoutes = make(map[string]Handler)
		}
		s.staticRoutes[path] = handler
	} else {
		if s.dynamicRoutes == nil {
			s.dynamicRoutes = make([]route, 0)
		}

		regex := "^" + path + "$"
		parts := strings.Split(path, "/")
		for _, part := range parts {
			if strings.HasPrefix(part, ":") {
				regex = strings.ReplaceAll(regex, part, fmt.Sprintf("(?P<%s>.*?)", part[1:]))
			}
		}

		s.dynamicRoutes = append(s.dynamicRoutes, route{
			regex:   regexp.MustCompile(regex),
			handler: handler,
		})
	}
}

// ListenAndServe starts the Server running on a specific port using the provided TLS configuration
func (s *Server) ListenAndServe(addr string, tlsConfig *tls.Config) error {
	// TODO: don't directly use tls.Config
	lInsecure, err := net.Listen("tcp", addr+":1965")
	if err != nil {
		return err
	}
	defer lInsecure.Close()

	l := tls.NewListener(lInsecure, tlsConfig)
	s.listener = l
	s.addr = addr

	defer l.Close()

	log.Info("Listening on ", "gemini://"+addr+":1965")
	for {
		conn, err := l.Accept()
		tlsConn, ok := conn.(*tls.Conn)
		if !ok {
			return errors.New("could not get tls connection")
		}
		if err != nil {
			return err
		}
		go s.handleConnection(tlsConn)
	}
}

func (s *Server) handleConnection(conn *tls.Conn) {
	defer conn.Close()

	request := make([]byte, 1026)
	_, err := conn.Read(request)
	if err != nil {
		log.Errorf("An error occurred while writing response: %s", err.Error())
		return
	}
	requestUri := strings.Split(string(request), "\r\n")[0]
	uri, err := url.Parse(requestUri)
	if err != nil {
		log.Error("Bad URI received")
		_, err := conn.Write([]byte("59 Bad Request\r\n"))
		if err != nil {
			log.Errorf("An error occurred while writing response: %s", err.Error())
		}
		return
	}

	if uri.Scheme != "gemini" {
		log.Error("Non-gemini URI received: " + requestUri)
		_, err := conn.Write([]byte("59 Only gemini URIs are supported (for now)\r\n"))
		if err != nil {
			log.Errorf("An error occurred while writing response: %s", err.Error())
		}
		return
	}

	handler, err := s.resolve(uri.Path)
	if err != nil {
		log.Error(uri.Path + " not found")
		_, err := conn.Write([]byte("51 Not Found\r\n"))
		if err != nil {
			log.Errorf("An error occurred while writing response: %s", err.Error())
		}
		return
	}

	handler(Request{
		URI:  *uri,
		conn: conn,
	})

	log.Info("Request received for " + strings.TrimRight(requestUri, "\r\n"))
}

func (s *Server) resolve(path string) (Handler, error) {
	handler, ok := s.staticRoutes[path]
	if ok {
		return handler, nil
	}

	for _, route := range s.dynamicRoutes {
		if route.regex.MatchString(path) {
			submatches := route.regex.FindStringSubmatch(path)
			params := make(map[string]string)
			for i, submatch := range submatches {
				if i == 0 {
					continue
				}

				name := route.regex.SubexpNames()[i]
				params[name] = submatch
			}

			return func(request Request) {
				request.Params = params
				route.handler(request)
			}, nil
		}
	}

	return nil, errors.New("route not found")
}
