package session

import "../models"

var SessionMain *Session

const (
	COOKIE_NAME = "sessionId"
)

type SessionData struct {
	Username string
}

type Session struct {
	Data map[string] *SessionData
}

func NewSession() *Session {
	s := new(Session)
	s.Data = make(map[string] *SessionData)
	return s
}

func (s *Session) Init (login string) string {
	id := models.GenerateId()
	data := &SessionData{login}
	s.Data[id] = data
	return id
}

func (s *Session) Get (id string) string {
	data := s.Data[id]
	if data == nil {
		return "";
	}
	return data.Username
}