package server

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"strings"
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
	write(conn, "Welcome!")
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		l := strings.TrimSpace(scanner.Text())
		input := strings.Split(l, " ")
		switch {
		case len(input) == 3 && input[0] == "SET":
			s.db.Set(input[1], input[2])
			write(conn, "KV-Pair added!")
		case len(input) == 2 && input[0] == "GET":
			if value, ok := s.db.Get(input[1]); ok {
				write(conn, value)
			} else {
				write(conn, "Key not found!")
			}
		case len(input) == 2 && input[0] == "DELETE":
			s.db.Delete(input[1])
			write(conn, "Deleted!")
		case len(input) == 1 && input[0] == "EXIT":
			conn.Close()
		default:
			write(conn, "Unknown Command!")
		}
		if _, ok := <-s.stop; !ok {
			write(conn, "Closing connection, server has stopped!")
			conn.Close()
			return
		}
	}
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

	//persisting to file
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(s.db)
	if err != nil {
		log.Fatal("GOB encoder error! ", err.Error())
	}
	err = ioutil.WriteFile("dbdata", buffer.Bytes(), 0600)
	if err != nil {
		log.Fatal("File write error! ", err.Error())
	}
	close(s.quit)
}

func write(conn net.Conn, s string) {
	_, err := fmt.Fprintln(conn, s)
	if err != nil {
		log.Fatal("Connection write error! ", err.Error())
	}
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
