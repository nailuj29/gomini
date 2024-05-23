package server

import (
	"crypto/tls"
	"errors"
	log "github.com/sirupsen/logrus"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

// A TitanHandler is a function to handle a [TitanRequest] by calling its various methods.
// The function is called when a titan request that it can handle is made as outlined in [Server.RegisterTitanHandler]
type TitanHandler func(request TitanRequest)

type titanRoute struct {
	regex   *regexp.Regexp
	handler TitanHandler
}

// RegisterTitanHandler sets up a [TitanHandler] to handle any [TitanRequest] that comes to a path
func (s *Server) RegisterTitanHandler(path string, handler TitanHandler) {
	if !strings.ContainsRune(path, ':') {
		if s.staticTitanRoutes == nil {
			s.staticTitanRoutes = make(map[string]TitanHandler)
		}
		s.staticTitanRoutes[path] = handler
	} else {
		if s.dynamicTitanRoutes == nil {
			s.dynamicTitanRoutes = make([]titanRoute, 0)
		}

		regex := createDynamicPathRegex(path)

		s.dynamicTitanRoutes = append(s.dynamicTitanRoutes, titanRoute{
			regex:   regexp.MustCompile(regex),
			handler: handler,
		})
	}
}

func (s *Server) handleTitanRequest(conn *tls.Conn, uri *url.URL) {
	rawParameters := strings.Split(uri.Path, ";")[1:]
	parameters := make(map[string]string)
	for _, rawParameter := range rawParameters {
		parameter := strings.Split(rawParameter, "=")
		if len(parameter) != 2 {
			log.Error("Malformed Parameter: " + rawParameter)
			_, err := conn.Write([]byte("59 Malformed parameter\r\n"))
			if err != nil {
				log.Errorf("An error occurred while writing response: %s", err.Error())
			}
			return
		}
		parameters[parameter[0]] = parameter[1]
	}

	titanRequest := TitanRequest{}
	token, ok := parameters["token"]
	if !ok {
		titanRequest.Token = ""
	} else {
		titanRequest.Token = token
	}

	mimeType, ok := parameters["mime"]
	if !ok {
		titanRequest.MIMEType = "text/gemini"
	} else {
		titanRequest.MIMEType = mimeType
	}

	size, ok := parameters["size"]
	if !ok {
		log.Error("Missing size parameter")
		_, err := conn.Write([]byte("59 Missing size parameter\r\n"))
		if err != nil {
			log.Errorf("An error occurred while writing response: %s", err.Error())
		}
		return
	} else {
		sizeInt, err := strconv.Atoi(size)
		if err != nil {
			log.Error("Malformed size parameter: " + size)
			_, err := conn.Write([]byte("59 Size must be a number\r\n"))
			if err != nil {
				log.Errorf("An error occurred while writing response: %s", err.Error())
			}
			return
		}
		body := make([]byte, sizeInt)
		_, err = conn.Read(body)
		if err != nil {
			log.Errorf("An error occurred while reading request body: %s", err.Error())
			return
		}

		titanRequest.Body = body

		handler, err := s.titanResolve(strings.Split(uri.Path, ";")[0])
		if err != nil {
			log.Error(uri.Path + " not found")
			_, err := conn.Write([]byte("51 Not Found\r\n"))
			if err != nil {
				log.Errorf("An error occurred while writing response: %s", err.Error())
			}

			return
		}

		titanRequest.conn = conn
		handler(titanRequest)
		log.Infof("Titan request received for %s", strings.TrimRight(uri.String(), "\r\n"))
	}
}

func (s *Server) titanResolve(path string) (TitanHandler, error) {
	handler, ok := s.staticTitanRoutes[path]
	if ok {
		return handler, nil
	}

	for _, route := range s.dynamicTitanRoutes {
		if route.regex.MatchString(path) {
			params := extractParams(path, route.regex)

			return func(request TitanRequest) {
				request.Params = params
				route.handler(request)
			}, nil
		}
	}

	return nil, errors.New("route not found")
}
