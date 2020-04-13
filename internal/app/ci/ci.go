package ci

import (
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/state"
	"github.com/RoboCup-SSL/ssl-game-controller/internal/app/tracker"
	"github.com/RoboCup-SSL/ssl-go-tools/pkg/sslconn"
	"log"
	"net"
	"time"
)

type Handler func(time.Time, *tracker.TrackerWrapperPacket) *state.Referee

type Server struct {
	address         string
	listener        net.Listener
	conn            net.Conn
	latestTime      time.Time
	tickChan        chan time.Time
	TrackerConsumer func(*tracker.TrackerWrapperPacket)
}

func NewServer(address string) *Server {
	return &Server{address: address, tickChan: make(chan time.Time, 1)}
}

func (s *Server) Start() {
	go s.listen(s.address)
}

func (s *Server) Stop() {
	if err := s.listener.Close(); err != nil {
		log.Printf("Could not close listener: %v", err)
	}
}

func (s *Server) listen(address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Print("Failed to listen on ", address)
		return
	}
	log.Print("Listening on ", address)
	s.listener = listener

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Could not accept connection: ", err)
		} else {
			log.Println("CI connection established")
			s.conn = conn
			s.serve(conn)
			log.Println("CI connection closed")
		}
	}
}

func (s *Server) serve(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("Could not close CI client connection: %v", err)
		}
	}()

	for {
		input := CiInput{}
		if err := sslconn.ReceiveMessage(conn, &input); err != nil {
			log.Println("Could not receive message from CI connection: ", err)
			return
		}

		if input.TrackerPacket != nil {
			s.TrackerConsumer(input.TrackerPacket)
		}

		if input.Timestamp != nil {
			sec := *input.Timestamp / 1e9
			nSec := *input.Timestamp - sec*1e9
			s.latestTime = time.Unix(sec, nSec)
			select {
			case s.tickChan <- time.Now():
			default:
			}
		}
	}
}

func (s *Server) SendMessage(refMsg *state.Referee) {
	if s.conn == nil {
		log.Println("Could not send message to CI client: Not connected")
		return
	}
	output := CiOutput{RefereeMsg: refMsg}
	if err := sslconn.SendMessage(s.conn, &output); err != nil {
		log.Printf("Could not send message: %v", err)
		return
	}
}

func (s *Server) Time() time.Time {
	return s.latestTime
}

func (s *Server) TickChanProvider() <-chan time.Time {
	return s.tickChan
}
