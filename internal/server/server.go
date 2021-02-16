package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/salmanahmed404/go-keyre/internal/store"
)

//Server is the multi-threaded TCP server
type Server struct {
	listener        net.Listener
	connections     map[uint]net.Conn
	stop            chan struct{}
	quit            chan struct{}
	connectionGroup *sync.WaitGroup
	db              *store.DB
}

func (s *Server) listen() {
	var connID uint
	connID = 1
	fmt.Println("Waiting for connections: ")
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if _, ok := <-s.stop; ok {
				log.Println("New connection error!", err.Error())
				continue
			} else {
				s.commit()
				break
			}
		}
		s.connections[connID] = conn
		s.connectionGroup.Add(1)
		go func(id uint) {
			defer s.connectionGroup.Done()
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
	close(s.stop)
	s.listener.Close()
	<-s.quit
}

func (s *Server) commit() {
	log.Println("Waiting for connections to close...")
	s.connectionGroup.Wait()
	log.Println("Commiting...")
	close(s.quit)
}

//NewServer creates a new instance of Server
func NewServer(service string) *Server {
	listener, err := net.Listen("tcp", service)
	if err != nil {
		log.Fatal("Listener Error! ", err.Error())
	}

	srv := &Server{
		listener:        listener,
		connections:     map[uint]net.Conn{},
		stop:            make(chan struct{}),
		quit:            make(chan struct{}),
		connectionGroup: new(sync.WaitGroup),
		db:              store.NewDB(),
	}
	go srv.listen()
	return srv
}
