package server

import (
	"fmt"
	"log"
	"net"
)

//Server is the multi-threaded TCP server
type Server struct {
	listener    net.Listener
	connections map[uint]net.Conn
	quit        chan struct{}
}

func (s *Server) listen() {
	var connID uint
	connID = 1
	fmt.Println("Waiting for connections: ")
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if _, ok := <-s.quit; ok {
				log.Println("New connection error!", err.Error())
				continue
			} else {
				s.commit()
				break
			}
		}
		s.connections[connID] = conn
		go func(id uint) {
			log.Printf("Connection with ID %d joined! ", id)
			s.handleConnection(conn)
			log.Printf("Connection with ID %d left! ", id)
			delete(s.connections, id)
		}(connID)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	conn.Write([]byte("Welcome!"))
	conn.Close()
}

//Stop is a...
func (s *Server) Stop() {
	log.Println("Server is stopping...")
	close(s.quit)
	s.listener.Close()
}

func (s *Server) commit() {
	log.Println("Commiting...")
}

//NewServer creates a new instance of Server
func NewServer(service string) *Server {
	listener, err := net.Listen("tcp", service)
	if err != nil {
		log.Fatal("Listener Error! ", err.Error())
	}

	srv := &Server{
		listener:    listener,
		connections: map[uint]net.Conn{},
		quit:        make(chan struct{}),
	}
	go srv.listen()
	return srv
}
