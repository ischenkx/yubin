package viewstat

import (
	"fmt"
	"net/http"
	"smtp-client/pkg/channel"
	"smtp-client/pkg/pubsub"
)

type Server struct {
	urlPrefix string
	handler   http.HandlerFunc
	pubsub    *pubsub.PubSub[Identifier]
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

func (s *Server) Visits() channel.Handle[Identifier] {
	return s.pubsub.Subscribe()
}
