package server

import (
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"
)

// A Handler is a function to handle a Request by calling its various methods.
// The function is called when a request that it can handle, as outlined in Server.RegisterHandler
type Handler func(request Request)

// A Server contains information required to run a TCP/TLS service capable of serving Gemini content over the internet
type Server struct {
	handlers map[string]Handler
	listener net.Listener
	addr     string
}

// New creates a new Server
func New() *Server {
	return &Server{}
}

// RegisterHandler sets up a Handler to handle any Request that comes to a path
func (s *Server) RegisterHandler(path string, handler Handler) {
	if s.handlers == nil {
		s.handlers = make(map[string]Handler)
	}
	s.handlers[path] = handler
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

	handler, ok := s.handlers[uri.Path]
	if !ok {
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
