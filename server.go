package session

import (
	"errors"
	"net/http"

	"github.com/mikerybka/util"
)

type Server struct {
	Handler  http.Handler
	Sessions map[string]bool
}

func getToken(r *http.Request) string {
	sessionCookie, err := r.Cookie("token")
	if err != nil {
		if errors.Is(err, http.ErrNoCookie) {
			return r.Header.Get("Token")
		} else {
			panic(err)
		}
	}
	return sessionCookie.Value
}

func (s *Server) newToken() string {
	token := util.RandomToken(64)
	for {
		if s.validSession(token) {
			token = util.RandomToken(64)
		} else {
			break
		}
	}
	return token
}

func (s *Server) validSession(token string) bool {
	return s.Sessions[token]
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token := getToken(r)
	if token == "" {
		token = s.newToken()
	} else if !s.validSession(token) {
		http.NotFound(w, r)
		return
	}
	r.Header.Set("SessionID", token)
	s.Handler.ServeHTTP(w, r)
}
