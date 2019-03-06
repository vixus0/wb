package wbd

import (
  "log"
  "time"
)

const (
  Sock = "/tmp/wbd.sock"
  LockDelay = 10 * time.Minute
)

type Server struct {
  session string
  stop *chan bool
}

func NewServer(sessionKey string, stop *chan bool) *Server {
  return &Server{session: sessionKey, stop: stop}
}

func (srv *Server) Stop(req int, resp *int) error {
  log.Println("stopping")
  *srv.stop <- true
  return nil
}

func (srv *Server) GetSession(req int, resp *string) error {
  log.Println("calling GetSession")
  *resp = srv.session
  return nil
}
