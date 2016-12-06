package socket

import (
	"log"
	"net"
	"os"
)

type Server struct {
	fd net.Listener
}

func NewServer() (*Server, error) {

	if Exists() {
		if err := os.Remove(FilePath()); err != nil {
			return nil, err
		}
	}

	fd, err := net.Listen("unix", FilePath())
	if err != nil {
		return nil, err
	}

	return &Server{fd: fd}, nil
}

func (s *Server) Run(clientHandle func(net.Conn)) {
	log.Printf("[INFO] Start listen %s\n", FilePath())
	for {
		clientFd, err := s.fd.Accept()
		if err != nil {
			log.Printf("[ERROR] %s\n", err.Error())
		}
		go clientHandle(clientFd)
	}
}

func (s *Server) Alive() bool {
	return Alive()
}

func (s *Server) Close() error {
	return s.fd.Close()
}
