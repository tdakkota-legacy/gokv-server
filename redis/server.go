package redis

import (
	"github.com/philippgille/gokv"
	"github.com/tidwall/redcon"
	"log"
	"strings"
)

type Server struct {
	store gokv.Store
	log   *log.Logger
}

func (s *Server) Handler(conn redcon.Conn, cmd redcon.Command) {
	command := strings.ToUpper(string(cmd.Args[0]))
	switch command {
	// system commands
	case "PING":
		if len(cmd.Args) > 2 {
			argFail(conn, command)
		} else if len(cmd.Args) == 2 {
			conn.WriteBulk(cmd.Args[1])
		} else {
			conn.WriteString("PONG")
		}
	case "QUIT":
		conn.WriteString("OK")
		if err := conn.Close(); err != nil {
			s.log.Println("error on closed", err)
		}

	// kv commands
	case "GET":
		if len(cmd.Args) != 2 {
			argFail(conn, command)
		} else {
			var val []byte
			if ok, err := s.store.Get(string(cmd.Args[1]), &val); !ok {
				conn.WriteNull()
				if err != nil {
					s.log.Println("get error caused:", err)
				}
			} else {
				conn.WriteBulk(val)
			}
		}
	case "SET":
		if len(cmd.Args) != 3 {
			argFail(conn, command)
		} else {
			if err := s.store.Set(string(cmd.Args[1]), cmd.Args[2]); err != nil {
				conn.WriteString("ERR driver error: " + err.Error())
				s.log.Println("set error caused:", err)
			} else {
				conn.WriteString("OK")
			}
		}
	case "DEL":
		if len(cmd.Args) < 2 {
			argFail(conn, command)
		} else {
			var n int
			for i := 1; i < len(cmd.Args); i++ {
				if err := s.store.Delete(string(cmd.Args[1])); err == nil {
					n++
				} else {
					s.log.Println("delete error caused:", err)
				}
			}
			conn.WriteInt(n)
		}
	default:
		conn.WriteError("ERR unknown command '" + command + "'")
	}
}

func ok(conn redcon.Conn) {
	conn.WriteString("OK")
}

func fail(conn redcon.Conn, s string) {
	conn.WriteError(s)
}

func argFail(conn redcon.Conn, command string) {
	fail(conn, "ERR wrong number of arguments for '"+command+"' command")
}

func (s *Server) Accept(conn redcon.Conn) bool {
	return true
}

func (s *Server) Closed(conn redcon.Conn, err error) {

}

func New(store gokv.Store, logger *log.Logger) *Server {
	return &Server{store: store, log: logger}
}
