package communication

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	"github.com/superc03/colirc/data"
)

type Server struct {
	listener net.Listener
	quit     chan interface{}
	wg       sync.WaitGroup
	l        *log.Logger
}

func NewServer(port int, l *log.Logger) *Server {
	s := &Server{
		quit: make(chan interface{}),
		l:    l,
	}
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.l.Fatal("Could Not Start Server")
	}
	s.listener = ln
	s.wg.Add(1)
	go s.serve()
	return s
}

func (s *Server) Stop() {
	close(s.quit)
	s.listener.Close()
	s.wg.Wait()
}

func (s *Server) serve() {
	defer s.wg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			select {
			case <-s.quit:
				return
			default:
				s.l.Println("TCP Accept Error")
			}
		} else {
			s.wg.Add(1)
			go func() {
				s.handleConnection(conn)
				s.wg.Done()
			}()
		}
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 2048)
	u := &data.Client{
		Hostname: conn.RemoteAddr().String(),
	}
ReadLoop:
	for {
		select {
		case <-s.quit:
			return
		default:
			conn.SetDeadline(time.Now().Add(200 * time.Millisecond))
			n, err := conn.Read(buf)
			if err != nil {
				if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
					continue ReadLoop
				} else if err != io.EOF {
					log.Println("Read Error", err)
					return
				}
			}
			if n == 0 {
				return
			}
			// Attempt to Convert to Message Struct
			msg, err := data.UnmarshalMessage(string(buf[:n-1]))
			if err != nil {
				BadRequestError(conn)
				continue ReadLoop
			}
			// Attempt to Use Function Handler, Otherwise Error Out
			handler := Handlers[msg.Command]
			if handler == nil {
				UnknownCommandError(conn, msg.Command)
				continue ReadLoop
			}
			handler(conn, u)
		}
	}
}
