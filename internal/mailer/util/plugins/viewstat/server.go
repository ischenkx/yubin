package viewstat

import (
	"fmt"
	"net/http"
	"sync"
)

type Server struct {
	mu        sync.RWMutex
	channels  []chan Identifier
	urlPrefix string
	handler   http.HandlerFunc
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.handler(writer, request)
}

func NewServer(urlPrefix string) *Server {
	return &Server{
		urlPrefix: urlPrefix,
	}
}

func (s *Server) GenerateLink(identifier Identifier) (string, error) {
	return fmt.Sprintf("%s/%s:%s.png", s.urlPrefix, identifier.User, identifier.Publication), nil
}

func (s *Server) Visits() Handle {
	s.mu.Lock()
	defer s.mu.Unlock()
	channel := make(chan Identifier, 2048)
	s.channels = append(s.channels, channel)
	return handle{
		channel: channel,
		server:  s,
	}
}

func (s *Server) publish(id Identifier) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.channels {
		select {
		case c <- id:
		default:
		}
	}
}

func (s *Server) removeChannel(c chan Identifier) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, candidate := range s.channels {
		if candidate == c {
			s.channels[i] = s.channels[len(s.channels)-1]
			s.channels = s.channels[:len(s.channels)-1]
			break
		}
	}
}

type handle struct {
	channel chan Identifier
	server  *Server
}

func (h handle) Chan() <-chan Identifier {
	return h.channel
}

func (h handle) Close() {
	h.server.removeChannel(h.channel)
}
