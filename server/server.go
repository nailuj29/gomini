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

// A Handler is a function to handle a [Request] by calling its various methods.
// The function is called when a request that it can handle is made, as outlined in [Server.RegisterHandler]
type Handler func(request Request)

// A Server contains information required to run a TCP/TLS service capable of serving Gemini content over the internet
type Server struct {
	staticRoutes       map[string]Handler
	staticTitanRoutes  map[string]TitanHandler
	dynamicRoutes      []route
	dynamicTitanRoutes []titanRoute
	listener           net.Listener
	addr               string
	running            bool
}

type route struct {
	regex   *regexp.Regexp
	handler Handler
}

// New creates a new [Server]
func New() *Server {
	return &Server{}
}

// RegisterHandler sets up a [Handler] to handle any [Request] that comes to a path
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

		regex := createDynamicPathRegex(path)

		s.dynamicRoutes = append(s.dynamicRoutes, route{
			regex:   regexp.MustCompile(regex),
			handler: handler,
		})
	}
}

func createDynamicPathRegex(path string) string {
	regex := "^" + path + "$"
	parts := strings.Split(path, "/")
	for _, part := range parts {
		if strings.HasPrefix(part, ":") {
			regex = strings.ReplaceAll(regex, part, fmt.Sprintf("(?P<%s>.*?)", part[1:]))
		}
	}
	return regex
}

// ListenAndServe starts the [Server] running on a specific port using the provided TLS configuration
func (s *Server) ListenAndServe(addr string, tlsConfig *tls.Config) error {
	// TODO: don't directly use tls.Config
	lInsecure, err := net.Listen("tcp", addr+":1965")
	if err != nil {
		return err
	}

	l := tls.NewListener(lInsecure, tlsConfig)
	s.listener = l
	s.addr = addr

	defer func() {
		if s.running {
			lInsecure.Close()
			l.Close()
		}
	}()

	log.Info("Listening on ", "gemini://"+addr+":1965")
	s.running = true
	for s.running {
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
	return nil
}

func (s *Server) Close() error {
	s.running = false
	return s.listener.Close()
}

func (s *Server) handleConnection(conn *tls.Conn) {
	defer conn.Close()

	request := make([]byte, 1026)
	buf := make([]byte, 1)
	var prevByte byte = 0
	i := 0
	done := false
	for !done {
		_, err := conn.Read(buf)
		if err != nil {
			log.Errorf("An error occurred while reading request %v", err)
			_, err := conn.Write([]byte("59 Bad Request\r\n"))
			if err != nil {
				log.Errorf("An error occurred while writing response: %s", err.Error())
			}
			return
		}
		request[i] = buf[0]
		i++
		if prevByte == 13 && buf[0] == 10 {
			done = true
		}
		prevByte = buf[0]
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

	if uri.Scheme != "gemini" && uri.Scheme != "titan" {
		log.Error("Non-gemini or titan URI received: " + requestUri)
		_, err := conn.Write([]byte("59 Only gemini and titan URIs are supported (for now)\r\n"))
		if err != nil {
			log.Errorf("An error occurred while writing response: %s", err.Error())
		}
		return
	}

	if uri.Scheme == "gemini" {
		s.handleGeminiRequest(conn, uri)
	} else {
		s.handleTitanRequest(conn, uri)
	}
}

func (s *Server) handleGeminiRequest(conn *tls.Conn, uri *url.URL) {
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

	log.Info("Gemini request received for " + strings.TrimRight(uri.String(), "\r\n"))
}

func (s *Server) resolve(path string) (Handler, error) {
	handler, ok := s.staticRoutes[path]
	if ok {
		return handler, nil
	}

	for _, route := range s.dynamicRoutes {
		if route.regex.MatchString(path) {
			params := extractParams(path, route.regex)

			return func(request Request) {
				request.Params = params
				route.handler(request)
			}, nil
		}
	}

	return nil, errors.New("route not found")
}

func extractParams(path string, regex *regexp.Regexp) map[string]string {
	submatches := regex.FindStringSubmatch(path)
	params := make(map[string]string)
	for i, submatch := range submatches {
		if i == 0 {
			continue
		}

		name := regex.SubexpNames()[i]
		params[name] = submatch
	}
	return params
}
